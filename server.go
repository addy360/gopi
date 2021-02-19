package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Panic(err)
	}
}
