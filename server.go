package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type postHandlers struct {
	store map[string]Post
}

type haha string

func main() {
	postHandler := newPost()
	http.HandleFunc("/posts", postHandler.get)

	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Panic(err)
	}

}

func (h *postHandlers) get(w http.ResponseWriter, r *http.Request) {
	posts := make([]Post, len(h.store))
	i := 0
	for _, v := range h.store {
		posts[i] = v
		i++
	}

	postBs, err := json.Marshal(posts)

	if err != nil {
		log.Panic(err)
	}

	w.Write(postBs)
}

func newPost() *postHandlers {
	return &postHandlers{
		store: map[string]Post{
			"postone": {
				ID:    "p1",
				Title: "title of post one",
				Body:  "body of post one",
			},
		},
	}
}
