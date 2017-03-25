# 4chan Data Dump

This is an attempt to start gathering text data dumps from various 4chan boards.
It uses the official [4chan-api](https://github.com/4chan/4chan-API).

The script stores everything `rethinkdb` for now.

### Table `threads`

Only few information for each thread is kept.
Additional information can be found from the first post of each thread that will
have an id of `board/threadno/threadno`.

```
{
    "id": "3/526260",
    "board": "3",
    "no": 526260,
    "archived": false,
    "last_modified": 1490396941
}
```

### Table `posts`

```
{
    "id": "3/248019/248019",
    "board": "3",
    "no": 248019,
    "archived_on": 0,
    "com": "So I was cleaning up my hard drive...",
    "ext": ".jpg",
    "filename": "3",
    "fsize": 3277,
    "h": 37,
    "images": 57,
    "md5": "zBbNPMGZwmwKcCb5AnW37A==",
    "name": "Anonymous",
    "now": "12/19/11(Mon)11:02",
    "replies": 73,
    "resto": 0,
    "semantic_url": "so-i-was-cleaning-up-my-hard-drive-and-i-found-my",
    "tim": 1324310530977,
    "time": 1324310530,
    "tn_h": 37,
    "tn_w": 145,
    "unique_ips": 11,
    "w": 145
}
```

### Table `entries`

This is a weird one; It takes "complex" posts such as ones that include 
replies to specific quotes and breaks them into multiple entries.

Take the following post for example:

```
Comment 1

>>100
Comment 2

>>200
>>300
>>400
Comment 3
```

This will be broken up into 5 entries:

* `Comment 1` in reply to the main thread.
* `Comment 2` in reply to post 100
* `Comment 3` in reply to post 200
* `Comment 3` in reply to post 300
* `Comment 3` in reply to post 400

Each of the entries will look something like this:

```
{
    "board": "3",
    "id": "3/248019/248019/0",
    "parent_id": 248019,
    "post_id": 248019,
    "text": "So I was cleaning up my hard drive..."
}
```