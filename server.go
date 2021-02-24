package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type postHandlers struct {
	sync.Mutex
	store map[string]Post
}

type haha string

func main() {
	postHandler := newPost()
	http.HandleFunc("/posts", postHandler.postController)
	http.HandleFunc("/post/", postHandler.getPost)

	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		log.Panic(err)
	}

}

func (h *postHandlers) postController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func (h *postHandlers) post(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Required body fields"))
		return
	}

	var post Post
	err = json.Unmarshal(bs, &post)
	h.Lock()
	post.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.store[post.ID] = post
	defer h.Unlock()
	w.Write([]byte("Success"))
}

func (h *postHandlers) get(w http.ResponseWriter, r *http.Request) {
	posts := make([]Post, len(h.store))
	i := 0
	h.Lock()
	for _, v := range h.store {
		posts[i] = v
		i++
	}
	h.Unlock()

	postBs, err := json.Marshal(posts)

	if err != nil {
		log.Panic(err)
	}

	w.Write(postBs)
}

func (h *postHandlers) getPost(w http.ResponseWriter, r *http.Request) {
	var post interface{}
	id := strings.Split(r.URL.Path, "/")[2]
	log.Println(id)
	i := 0
	h.Lock()
	for _, v := range h.store {
		if v.ID == strings.TrimSpace(id) {
			post = v
		}
		i++
	}
	h.Unlock()

	postBs, err := json.Marshal(post)

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
