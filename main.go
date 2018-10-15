package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cpachomski/hacker-news-api/router"
)

func main() {

	r := router.CreateRouter()

	fmt.Println("Serving up something good on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
