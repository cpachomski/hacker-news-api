package main

type Story struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	By           string `json:"by"`
	Score        int    `json:"score"`
	Time         int    `json:"time"`
	Comments     []int  `json:"comments"`
	CommentCount int    `json:"commentsCount"`
}

type Stories []Story
