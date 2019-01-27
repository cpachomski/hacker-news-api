package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Story struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	By           string `json:"by"`
	Score        int    `json:"score"`
	Time         int    `json:"time"`
	Comments     []int  `json:"kids"`
	CommentCount int    `json:"commentsCount"`
}

type Stories []Story

func fetch(url string, ch chan<- []byte) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	ch <- body
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// List Type: top | best | new | ask | show | job
	listType := "top"
	if _, ok := request.PathParameters["type"]; ok {
		listType = request.PathParameters["type"]
	}

	// Fetch list of story ids
	topStoriesIdsCh := make(chan []byte)
	topStoriesURI := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/%sstories.json", listType)
	go fetch(topStoriesURI, topStoriesIdsCh)

	var topStoriesIds []int
	err := json.Unmarshal(<-topStoriesIdsCh, &topStoriesIds)
	if err != nil {
		panic(err)
	}

	limit := 30
	if _, ok := request.QueryStringParameters["limit"]; ok {
		limit, err = strconv.Atoi(request.QueryStringParameters["limit"])
		if err != nil {
			panic(err)
		}
	}

	topStoriesRange := topStoriesIds[0:limit]
	topStoriesCh := make(chan []byte)
	for _, storyId := range topStoriesRange {
		storyURI := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyId)
		go fetch(storyURI, topStoriesCh)
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

	result, err := json.Marshal(topStories)
	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{Body: string(result), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
