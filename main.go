package main

import (
	"log"
	"net/http"
	"time"
)

// store is the global session store
var store *SessionStore

func main() {
	store = NewSessionStore(5 * time.Minute)
	go store.Cleanup()

	// Route handlers
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/respond", handleRespond)
	// Redirect root to /start
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/start", http.StatusFound)
	})

	log.Println("Server starting on 127.0.0.1:8080 (run as root for pam_chauthtok)")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
