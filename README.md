# RSSer

[![GoDoc](https://godoc.org/github.com/benjaminheng/rsser?status.svg)](https://godoc.org/github.com/benjaminheng/rsser)

RSSer allows you to generate RSS feeds for services that don't normally provide
RSS feeds.

Currently RSSer supports the following sites:

- Instagram

Contributions for additional sites is strongly welcome!

## Usage

Here's a minimal service that exposes RSS feeds for Instagram users.

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/benjaminheng/rsser/instagram"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/feed/instagram/user/{username}", InstagramFeedHandler)
	log.Println("Server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func InstagramFeedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	feed, err := instagram.GetUserFeed(context.Background(), username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	rss, err := feed.ToRss()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rss))
}
```

See also: [cmd/server/main.go](./cmd/server/main.go).
