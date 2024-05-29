package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	urlMap = make(map[string]string)
	mu     sync.RWMutex
)

func main() {
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/", redirectHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

type Url struct {
	OrgnalUrl string `json:"orignal_url"`
}

type ShortenedUrl struct {
	ShortenUrl string `json:"shorten_url"`
}

type Error struct {
	StatusCode int
	Message    string `json:""`
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	url := Url{}
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		formatedError := Error{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
		json.NewEncoder(w).Encode(formatedError)
	}

	b := make([]byte, 6)
	rand.Read(b)
	shortenUrl := strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
	fmt.Println(shortenUrl)
	mu.Lock()
	urlMap[shortenUrl] = url.OrgnalUrl
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ShortenedUrl{
		ShortenUrl: shortenUrl,
	})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortened := r.URL.Path[1:]

	mu.RLock()
	original, ok := urlMap[shortened]
	mu.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}
