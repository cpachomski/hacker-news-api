package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/creack/httpreq"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Todos API v0.0.1"))
}

func TopStories(w http.ResponseWriter, r *http.Request) {
	// parse query params
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := &QueryParams{}
	defaultLimit := 30

	if err := (httpreq.ParsingMap{
		{Field: "limit", Fct: httpreq.ToInt, Dest: &params.Limit},
		{Field: "page", Fct: httpreq.ToInt, Dest: &params.Page},
	}.Parse(r.Form)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// GET list of top 500 stories
	topStoriesIdsCh := make(chan []byte)
	topStoriesURI := "https://hacker-news.firebaseio.com/v0/topstories.json"
	go MakeRequest(topStoriesURI, topStoriesIdsCh)

	var topStoriesIds []int
	err := json.Unmarshal(<-topStoriesIdsCh, &topStoriesIds)
	if err != nil {
		panic(err)
	}

	var topStoriesRange []int
	if params.Page < 2 {
		if params.Limit > 0 {
			topStoriesRange = topStoriesIds[0:params.Limit]
		} else {
			topStoriesRange = topStoriesIds[0:defaultLimit]
		}
	} else {
		if params.Limit > 0 {
			topStoriesRange = topStoriesIds[params.Page*params.Limit : (params.Page*params.Limit)+params.Limit]
		} else {
			topStoriesRange = topStoriesIds[params.Page*defaultLimit : (params.Page*defaultLimit)+defaultLimit]
		}
	}

	topStoriesCh := make(chan []byte)

	for _, storyId := range topStoriesRange {
		storyURI := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyId)
		go MakeRequest(storyURI, topStoriesCh)
	}

	var topStories Stories
	for range topStoriesRange {
		var story Story
		err := json.Unmarshal(<-topStoriesCh, &story)
		if err != nil {
			panic(err)
		}

		topStories = append(topStories, story)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(topStories); err != nil {
		panic(err)
	}
}
