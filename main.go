package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	router := NewRouter()

	fmt.Println("Serving up something good on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
