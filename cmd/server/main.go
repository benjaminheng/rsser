package main

import (
	"context"
	"log"
	"net/http"

	"github.com/benjaminheng/rsser/rsser/instagram"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/feed/instagram/user/{username}", InstagramFeedHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}

func InstagramFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	username := vars["username"]
	feed, err := instagram.GetUserFeed(ctx, username)
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
