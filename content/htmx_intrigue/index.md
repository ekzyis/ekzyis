---
title: HTMX Intrigue
date: 2024-06-07T06:48:38.417Z
sn_id: 564849
---

> ### htmx in a nutshell
>
> htmx is a library that allows you to access modern browser features directly from HTML, rather than using javascript.

-- [htmx.org/docs/#introduction](https://htmx.org/docs/#introduction)

A few days ago and after having read a lot about [HTMX](https://htmx.org/) and [HATEOAS](https://htmx.org/essays/hateoas/), I finally got my feet wet with HTMX.

I wanted to include the amount of comments and sats that I received on SN on my [static site](https://ekzy.is/about/) that is generated with [Hugo](https://gohugo.io/):

![](https://m.stacker.news/34338)

The cool thing about HTMX is that you don't need to write javascript yourself. The HTML spec is extended with attributes like `hx-get`, `hx-trigger`, `hx-target` [and more](https://htmx.org/reference/) which declare the behavior you want.

For example, this is the HTML fragment that I included in my [HTML template](https://git.ekzy.is/ekzyis/ekzyis/src/commit/6fc1eb2f7d59448d317414becb0a6a0721707b07/layouts/_default/single.html):

```html
{{ with .Params.sn_id }}
<span
    hx-get="/api/content_meta?sn_id={{- . -}}"
    hx-trigger="load"
    hx-swap="outerHTML"
>
    <span class="px-1">|</span>
    <span>
        <a
            class="underline"
            href="https://stacker.news/items/{{- . -}}"
            target="_blank"
            rel="noopener noreferrer me">0 comments</a>
    </span>
    <span class="px-1">|</span>
    <span>0 sats</span>
<span>
{{ end }}
```

On page load (`hx-trigger="load"`), a GET request to /api/content_meta?sn_id=<id> is triggered (`hx-get`) and the response replaces `<span>` (`hx-target`).[^1] 

For the backend, I went with go and [`net/http`](https://pkg.go.dev/net/http) which is part of go's stdlib. I also learned about [`templ`](https://templ.guide/) which makes it possible to write something like JSX in go.[^2]

This is how I wrote the HTML:

```go
package main

import "github.com/ekzyis/snappy"
import "fmt"

templ contentMeta (item *sn.Item) {
    <span class="px-1">|</span>
    <span>
        <a
            class="underline"
            href={ templ.URL(fmt.Sprintf("https://stacker.news/items/%d", item.Id)) }
            target="_blank"
            rel="noopener noreferrer me"
        >
        { fmt.Sprintf("%d comments", item.NComments) }
        </a>
    </span>
    <span class="px-1">|</span>
    <span>{ fmt.Sprintf("%d sats", item.Sats) }</span>
}
```

And then I ran `templ generate` in the terminal which generated a go file which includes a function that I can call in my [request handler](https://git.ekzy.is/ekzyis/ekzyis/src/commit/6fc1eb2f7d59448d317414becb0a6a0721707b07/server/server.go) to render the HTML:

```go
func HandleContentMeta(w http.ResponseWriter, r *http.Request) {
	var (
		// stacker.news item id
		qid  = r.URL.Query().Get("sn_id")
		c    = sn.NewClient()
		item *sn.Item
		id   int
		err  error
	)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		return
	}

	if qid == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "sn_id query param required\n")
		return
	}

	if id, err = strconv.Atoi(qid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "sn_id must be numeric\n")
		return
	}

	if item, err = c.Item(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	contentMeta(item).Render(r.Context(), w)
}
```

_The heavy lifting of fetching the SN item data is done with [snappy](https://github.com/ekzyis/snappy), the library that I wrote for [@hn](https://stacker.news/hn) and [@nitter](https://stacker.news/nitter)._

And that is all! Quite simple, ain't it?

I am really intrigued by HTMX now, even more now than I already was before trying it out. I suspected it wouldn't be that easy in practice. Of course, this was only very basic usage but at least this very basic usage is indeed very basic!

But now I wonder ... are there any HTMX devs here? If so, I'd be interested in your experience. Do you like it? What did you build with it? What do you not like about it?

[^1]: For example, [this](https://ekzy.is/api/content_meta?sn_id=564849) is the response for this post.

[^2]: JSX was my original mindblow of [React](https://react.dev/)