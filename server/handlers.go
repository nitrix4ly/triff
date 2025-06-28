package server

import (
	"encoding/json"
	"net/http"
	"github.com/nitrix4ly/triff/core"
)

type Handler struct {
	DB *core.Database
}

func NewHandler(db *core.Database) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, err := h.DB.Get(key)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"value": value})
}

func (h *Handler) SetHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	key := body["key"]
	value := body["value"]
	err := h.DB.Set(key, value)
	if err != nil {
		http.Error(w, "Failed to set key", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	err := h.DB.Delete(key)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
