package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/benjaminheng/rsser/rsser/feedcache"
	"github.com/benjaminheng/rsser/rsser/instagram"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
)

type server struct {
	cache    *feedcache.Cache
	cacheTTL time.Duration
}

func main() {
	ttl, err := time.ParseDuration("15m")
	if err != nil {
		log.Fatal(err)
	}
	srv := &server{
		cache:    feedcache.New(),
		cacheTTL: ttl,
	}
	r := mux.NewRouter()
	r.HandleFunc("/feed/instagram/user/{username}", srv.InstagramFeedHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}

func (s *server) InstagramFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	username := vars["username"]
	key := fmt.Sprintf("instagram:%s", username)
	if feed := s.cache.Get(key); feed != nil {
		writeFeedToResponse(w, feed)
		return
	}
	feed, err := instagram.GetUserFeed(ctx, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	s.cache.Set(key, feed, s.cacheTTL)
	writeFeedToResponse(w, feed)
}

func writeFeedToResponse(w http.ResponseWriter, feed *feeds.Feed) {
	rss, err := feed.ToRss()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rss))
}
