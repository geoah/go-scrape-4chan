package main

import (
	"fmt"
)

// Page -
type Page struct {
	Page    int       `json:"page"`
	Threads []*Thread `json:"threads"`
}

// Thread -
type Thread struct {
	ID           string `json:"id" gorethink:"id,omitempty"`
	Board        string `json:"board" gorethink:"board"`
	No           int    `json:"no" gorethink:"no"`
	LastModified int    `json:"last_modified" gorethink:"last_modified"`
	Archived     bool   `json:"archived" gorethink:"archived"`
}

// Posts -
type Posts struct {
	Posts []*Post `json:"posts"`
}

// Boards -
type Boards struct {
	Boards []*Board `json:"boards"`
}

// Board -
type Board struct {
	Board           string `json:"board"`
	Title           string `json:"title"`
	WsBoard         int    `json:"ws_board"`
	PerPage         int    `json:"per_page"`
	Pages           int    `json:"pages"`
	MaxFilesize     int    `json:"max_filesize"`
	MaxWebmFilesize int    `json:"max_webm_filesize"`
	MaxCommentChars int    `json:"max_comment_chars"`
	MaxWebmDuration int    `json:"max_webm_duration"`
	BumpLimit       int    `json:"bump_limit"`
	ImageLimit      int    `json:"image_limit"`
	Cooldowns       struct {
		Threads int `json:"threads"`
		Replies int `json:"replies"`
		Images  int `json:"images"`
	} `json:"cooldowns"`
	MetaDescription string `json:"meta_description"`
	IsArchived      int    `json:"is_archived,omitempty"`
	Spoilers        int    `json:"spoilers,omitempty"`
	CustomSpoilers  int    `json:"custom_spoilers,omitempty"`
	ForcedAnon      int    `json:"forced_anon,omitempty"`
	UserIds         int    `json:"user_ids,omitempty"`
	CodeTags        int    `json:"code_tags,omitempty"`
	WebmAudio       int    `json:"webm_audio,omitempty"`
	MinImageWidth   int    `json:"min_image_width,omitempty"`
	MinImageHeight  int    `json:"min_image_height,omitempty"`
	Oekaki          int    `json:"oekaki,omitempty"`
	CountryFlags    int    `json:"country_flags,omitempty"`
	SjisTags        int    `json:"sjis_tags,omitempty"`
	TextOnly        int    `json:"text_only,omitempty"`
	RequireSubject  int    `json:"require_subject,omitempty"`
	MathTags        int    `json:"math_tags,omitempty"`
}

// Post -
type Post struct {
	ID          string `json:"id" gorethink:"id,omitempty"`
	Board       string `json:"board" gorethink:"board"`
	No          int    `json:"no" gorethink:"no"`
	Now         string `json:"now" gorethink:"now"`
	Name        string `json:"name" gorethink:"name"`
	Com         string `json:"com,omitempty" gorethink:"com,omitempty"`
	Filename    string `json:"filename,omitempty" gorethink:"filename,omitempty"`
	Ext         string `json:"ext,omitempty" gorethink:"ext,omitempty"`
	W           int    `json:"w,omitempty" gorethink:"w,omitempty"`
	H           int    `json:"h,omitempty" gorethink:"h,omitempty"`
	TnW         int    `json:"tn_w,omitempty" gorethink:"tn_w,omitempty"`
	TnH         int    `json:"tn_h,omitempty" gorethink:"tn_h,omitempty"`
	Tim         int64  `json:"tim,omitempty" gorethink:"tim,omitempty"`
	Time        int    `json:"time" gorethink:"time"`
	Md5         string `json:"md5,omitempty" gorethink:"md5,omitempty"`
	Fsize       int    `json:"fsize,omitempty" gorethink:"fsize,omitempty"`
	Resto       int    `json:"resto" gorethink:"resto"`
	Bumplimit   int    `json:"bumplimit,omitempty" gorethink:"bumplimit,omitempty"`
	Imagelimit  int    `json:"imagelimit,omitempty" gorethink:"imagelimit,omitempty"`
	SemanticURL string `json:"semantic_url,omitempty" gorethink:"semantic_url,omitempty"`
	Replies     int    `json:"replies,omitempty" gorethink:"replies,omitempty"`
	Images      int    `json:"images,omitempty" gorethink:"images,omitempty"`
	UniqueIps   int    `json:"unique_ips,omitempty" gorethink:"unique_ips,omitempty"`
	TailSize    int    `json:"tail_size,omitempty" gorethink:"tail_size,omitempty"`
	ArchivedOn  int    `json:"archived_on" gorethink:"archived_on"`
}

// NewPost -
func NewPost(board string, threadNo int, post *Post) *Post {
	post.ID = fmt.Sprintf("%s/%d/%d", board, threadNo, post.No)
	post.Board = board
	return post
}

// Entry -
type Entry struct {
	ID       string `gorethink:"id,omitempty"`
	Board    string `gorethink:"board"`
	PostID   int    `gorethink:"post_id"`
	ParentID int    `gorethink:"parent_id"`
	Text     string `gorethink:"text"`
}

// NewEntry -
func NewEntry(board string, postID, parentID, i int, text string) *Entry {
	return &Entry{
		ID:       fmt.Sprintf("%s/%d/%d/%d", board, postID, parentID, i),
		Board:    board,
		PostID:   postID,
		ParentID: parentID,
		Text:     text,
	}
}
