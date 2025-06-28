package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nitrix4ly/triff/core"
	"github.com/nitrix4ly/triff/utils"
)

var db *core.Database

func main() {
	// Initialize DB
	db = core.NewDatabase()

	port := getPort()
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Printf("Web server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}
	val, err := db.Get(key)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"value": val})
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	key, ok1 := payload["key"]
	value, ok2 := payload["value"]
	if !ok1 || !ok2 {
		http.Error(w, "Missing key or value", http.StatusBadRequest)
		return
	}
	db.Set(key, value)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}
	if err := db.Delete(key); err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
