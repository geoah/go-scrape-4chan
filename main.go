package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"os"

	r "gopkg.in/gorethink/gorethink.v3"
)

const (
	boardsURL  = "http://a.4cdn.org/boards.json"
	archiveURL = "http://a.4cdn.org/%s/archive.json"
	threadsURL = "http://a.4cdn.org/%s/threads.json"
	threadURL  = "http://a.4cdn.org/%s/thread/%d.json"

	threadsTable = "threads"
	postsTable   = "posts"
	entriesTable = "entries"
)

var (
	myClient = &http.Client{Timeout: 10 * time.Second}
	tagsRe   = regexp.MustCompile(`<[^>]*>`)
	spacesRe = regexp.MustCompile(`\s\s+`)
	quotesRe = regexp.MustCompile(`>>([0-9]+) `)
)

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost",
		Database: "4chan",
	})
	if err != nil {
		log.Fatalln(err)
	}

	persistence := &RethinkPersistence{session}

	boards, err := getBoards()
	if err != nil {
		log.Println("*** Could not get boards", err)
		os.Exit(1)
	}
	for ib, board := range boards {
		fmt.Printf("> Going through /%s/ (%d out of %d)\n", board.Board, ib+1, len(boards))
		threads, err := getThreads(board.Board, true)
		if err != nil {
			log.Println("*** Could not get threads:", err)
			continue
		}
		for it, thread := range threads {
			thread.ID = fmt.Sprintf("%s/%d", board.Board, thread.No)
			thread.Board = board.Board

			fmt.Printf(">> Going through thread %s (archived: %t) (%d out of %d)\n", thread.ID, thread.Archived, it+1, len(threads))

			if oThread, err := persistence.GetThread(thread.ID); err != nil {
				if err != r.ErrEmptyResult {
					fmt.Println("*** Could not get thread:", err)
				}
			} else {
				if oThread.Archived {
					fmt.Println(">>> Skipping thread (archived and already indexed):", thread.ID)
					continue
				}
				if thread.Archived == false {
					if oThread.LastModified >= thread.LastModified {
						fmt.Println(">>> Skipping thread (not modified):", thread.ID)
						continue
					}
				}
			}
			posts, err := getPosts(board.Board, thread.No)
			if err != nil {
				log.Println("*** Could not get posts", err)
				continue
			}
			if thread.Archived {
				thread.LastModified = posts[0].ArchivedOn
			}
			fmt.Printf(">>> ")
			for _, post := range posts {
				post := NewPost(board.Board, thread.No, post)
				if err := persistence.Put(postsTable, true, post); err != nil {
					if r.IsConflictErr(err) {
						fmt.Printf("=")
						continue
					}
					log.Println("*** Could not write post:", err)
				}
				fmt.Printf("+")
				text := cleanupHTML(post.Com)
				if len(text) > 0 {
					quotes, texts := splitMessage(strconv.Itoa(thread.No), text)
					for i, quote := range quotes {
						quoteNo, _ := strconv.Atoi(quote)
						entry := NewEntry(board.Board, post.No, quoteNo, i, texts[i])
						if err := persistence.Put(entriesTable, false, entry); err != nil {
							log.Println("*** Could not write entry:", err)
						}
					}
				}
			}
			fmt.Println()
			// once we are done processing all posts for this tread
			// we can store the thread with its last updated so we don't have to
			// reprocess if nothing has changed
			if err := persistence.Put(threadsTable, false, thread); err != nil {
				log.Println("*** Could not put thread:", err)
			}
		}
	}
}

func splitMessage(ono, otext string) ([]string, []string) {
	otext = strings.TrimLeft(otext, " ") + " "
	ms := quotesRe.FindAllStringSubmatchIndex(otext, -1)
	squotes := []string{}
	texts := []string{}
	for i, m := range ms {
		if i == 0 && m[i] > 0 {
			squotes = append(squotes, "")
			texts = append(texts, strings.Trim(otext[:ms[i][0]], " "))
		}

		rt := ""
		ct := ""
		rt = otext[m[2]:m[3]]
		if i+1 < len(ms) {
			ct = otext[m[1]:ms[i+1][0]]
		} else {
			ct = otext[m[1]:]
		}
		squotes = append(squotes, rt)
		texts = append(texts, strings.Trim(ct, " "))
	}

	fill := []int{}
	for i, ct := range texts {
		if len(ct) == 0 {
			fill = append(fill, i)
		} else {
			for _, j := range fill {
				texts[j] = ct
			}
			fill = []int{}
		}
	}
	for _, j := range fill {
		texts[j] = texts[len(texts)-1]
	}

	for j := range squotes {
		if squotes[j] == "" {
			squotes[j] = ono
		}
	}

	if len(squotes) != len(texts) {
		log.Fatal("Quotes should always be the same as texts.")
	}

	if len(squotes) == 0 {
		squotes = append(squotes, ono)
		texts = append(texts, otext)
	}

	return squotes, texts
}

func cleanupHTML(text string) string {
	text = html.UnescapeString(text)
	text = tagsRe.ReplaceAllString(text, " ")
	text = strings.Replace(text, "\n", " ", -1)
	text = strings.Replace(text, "\r", " ", -1)
	text = strings.Replace(text, "\t", " ", -1)
	text = strings.Trim(text, " ")
	text = spacesRe.ReplaceAllString(text, " ")
	return text
}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if r != nil && r.StatusCode != 200 {
		fmt.Println("*** Bad status code:", r.StatusCode)
	}
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func getBoards() ([]*Board, error) {
	boards := &Boards{}
	if err := getJSON(boardsURL, &boards); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return boards.Boards, nil
}

func getThreads(boardTag string, includeArchive bool) ([]*Thread, error) {
	board := []*Page{}
	threads := []*Thread{}
	url := fmt.Sprintf(threadsURL, boardTag)
	if err := getJSON(url, &board); err != nil {
		fmt.Println("Error:", err)
		return threads, err
	}
	for _, page := range board {
		for _, thread := range page.Threads {
			threads = append(threads, thread)
		}
	}
	if includeArchive {
		archivedThreads := []int{}
		url := fmt.Sprintf(archiveURL, boardTag)
		if err := getJSON(url, &archivedThreads); err != nil {
			fmt.Println("Error archive:", err)
			return threads, err
		}
		for _, threadNo := range archivedThreads {
			threads = append(threads, &Thread{
				No:       threadNo,
				Archived: true,
			})
		}
	}
	return threads, nil
}

func getPosts(boardTag string, threadNo int) ([]*Post, error) {
	posts := &Posts{}
	url := fmt.Sprintf(threadURL, boardTag, threadNo)
	if err := getJSON(url, &posts); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return posts.Posts, nil
}
