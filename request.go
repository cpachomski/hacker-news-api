package main

import (
	"io/ioutil"
	"net/http"
)

func MakeRequest(url string, ch chan<- []byte) {
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
