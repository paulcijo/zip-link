package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/paulcijo/zip-link/redis"
)

const (
	counterKey = "urlCounts"
)

// createKey creates a key based on an int Id for the URL
func createKey(_id int) string {
	_idStr := strconv.Itoa(_id)
	data := base64.StdEncoding.EncodeToString([]byte(_idStr))
	return data
}

func shortURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	longURL, err := redis.Get(vars["shortURL"])
	if err != nil {
		w.Write([]byte("Invalid URL"))
		return
	}
	http.Redirect(w, r, string(longURL), http.StatusPermanentRedirect)
}

func longURLHandler(w http.ResponseWriter, r *http.Request) {
	url := []byte(r.URL.Path[5:])
	count, err := redis.Incr(counterKey)
	if err != nil {
		panic(err)
	}
	key := createKey(count)
	redis.Set(key, url)
	w.Write([]byte(key))
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/new/").HandlerFunc(longURLHandler)
	r.HandleFunc("/{shortURL}", shortURLHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Ping the connection to check if it works
	log.Fatal(srv.ListenAndServe())
}
