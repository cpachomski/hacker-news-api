package main

import (
	"net/http"

	"github.com/cpachomski/hacker-news-api/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"TopStories",
		"GET",
		"/top-stories",
		handlers.TopStories,
	},
}
