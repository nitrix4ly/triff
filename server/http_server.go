package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nitrix4ly/triff/commands"
	"github.com/nitrix4ly/triff/core"
	"github.com/nitrix4ly/triff/utils"
)

// HTTPServer handles HTTP REST API requests
type HTTPServer struct {
	db             *core.Database
	port           int
	router         *mux.Router
	stringCommands *commands.StringCommands
	logger         *utils.Logger
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(db *core.Database, port int, logger *utils.Logger) *HTTPServer {
	server := &HTTPServer{
		db:             db,
		port:           port,
		router:         mux.NewRouter(),
		stringCommands: commands.NewStringCommands(db),
		logger:         logger,
	}
	
	server.setupRoutes()
	return server
}

// Start begins the HTTP server
func (s *HTTPServer) Start() error {
	s.logger.Info(fmt.Sprintf("HTTP server listening on port %d", s.port))
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}

// setupRoutes configures all HTTP routes
func (s *HTTPServer) setupRoutes() {
	// Add CORS middleware
	s.router.Use(s.corsMiddleware)
	s.router.Use(s.loggingMiddleware)

	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	
	// Basic operations
	api.HandleFunc("/ping", s.handlePing).Methods("GET")
	api.HandleFunc("/info", s.handleInfo).Methods("GET")
	api.HandleFunc("/keys", s.handleKeys).Methods("GET")
	api.HandleFunc("/keys/{key}", s.handleKeyOperations).Methods("GET", "POST", "PUT", "DELETE")
	api.HandleFunc("/keys/{key}/ttl", s.handleTTL).Methods("GET", "POST")
	api.HandleFunc("/keys/{key}/exists", s.handleExists).Methods("GET")
	
	// String operations
	api.HandleFunc("/string/{key}", s.handleStringGet).Methods("GET")
	api.HandleFunc("/string/{key}", s.handleStringSet).Methods("POST", "PUT")
	api.HandleFunc("/string/{key}/append", s.handleStringAppend).Methods("POST")
	api.HandleFunc("/string/{key}/length", s.handleStringLength).Methods("GET")
	api.HandleFunc("/string/{key}/incr", s.handleStringIncr).Methods("POST")
	api.HandleFunc("/string/{key}/decr", s.handleStringDecr).Methods("POST")
	
	// Bulk operations
	api.HandleFunc("/bulk/get", s.handleBulkGet).Methods("POST")
	api.HandleFunc("/bulk/set", s.handleBulkSet).Methods("POST")
	api.HandleFunc("/flush", s.handleFlushAll).Methods("DELETE")
}

// Middleware functions
func (s *HTTPServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}

// Handler functions
func (s *HTTPServer) handlePing(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "PONG", "status": "ok"}
	s.writeJSON(w, http.StatusOK, response)
}

func (s *HTTPServer) handleInfo(w http.ResponseWriter, r *http.Request) {
	info := s.db.Info()
	s.writeJSON(w, http.StatusOK, info)
}

func (s *HTTPServer) handleKeys(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")
	if pattern == "" {
		pattern = "*"
	}
	
	keys := s.db.Keys(pattern)
	response := map[string]interface{}{
		"keys":  keys,
		"count": len(keys),
	}
	s.writeJSON(w, http.StatusOK, response)
}

func (s *HTTPServer) handleKeyOperations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	switch r.Method {
	case "GET":
		value, exists := s.db.Get(key)
		if !exists {
			s.writeError(w, http.StatusNotFound, "key not found")
			return
		}
		
		response := map[string]interface{}{
			"key":        key,
			"value":      value.Data,
			"type":       value.Type,
			"ttl":        s.db.GetTTL(key),
			"created_at": value.CreatedAt,
			"updated_at": value.UpdatedAt,
		}
		s.writeJSON(w, http.StatusOK, response)
		
	case "DELETE":
		if s.db.Delete(key) {
			s.writeJSON(w, http.StatusOK, map[string]string{"message": "key deleted"})
		} else {
			s.writeError(w, http.StatusNotFound, "key not found")
		}
	}
}

func (s *HTTPServer) handleStringGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	response := s.stringCommands.Get(key)
	if response.Success {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"key":   key,
			"value": response.Data,
		})
	} else {
		s.writeError(w, http.StatusNotFound, "key not found or not a string")
	}
}

func (s *HTTPServer) handleStringSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	var payload struct {
		Value string `json:"value"`
		TTL   int64  `json:"ttl,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	
	response := s.stringCommands.Set(key, payload.Value, payload.TTL)
	if response.Success {
		s.writeJSON(w, http.StatusOK, map[string]string{"message": "value set successfully"})
	} else {
		s.writeError(w, http.StatusInternalServerError, response.Error)
	}
}

func (s *HTTPServer) handleStringAppend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	var payload struct {
		Value string `json:"value"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	
	response := s.stringCommands.Append(key, payload.Value)
	if response.Success {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "value appended successfully",
			"length":  response.Data,
		})
	} else {
		s.writeError(w, http.StatusInternalServerError, response.Error)
	}
}

func (s *HTTPServer) handleStringLength(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	response := s.stringCommands.Strlen(key)
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"key":    key,
		"length": response.Data,
	})
}

func (s *HTTPServer) handleStringIncr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	var payload struct {
		By int64 `json:"by,omitempty"`
	}
	
	json.NewDecoder(r.Body).Decode(&payload)
	
	var response *core.Response
	if payload.By == 0 {
		response = s.stringCommands.Incr(key)
	} else {
		response = s.stringCommands.IncrBy(key, payload.By)
	}
	
	if response.Success {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"key":   key,
			"value": response.Data,
		})
	} else {
		s.writeError(w, http.StatusBadRequest, response.Error)
	}
}

func (s *HTTPServer) handleStringDecr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	response := s.stringCommands.Decr(key)
	if response.Success {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"key":   key,
			"value": response.Data,
		})
	} else {
		s.writeError(w, http.StatusBadRequest, response.Error)
	}
}

func (s *HTTPServer) handleTTL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	switch r.Method {
	case "GET":
		ttl := s.db.GetTTL(key)
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"key": key,
			"ttl": ttl,
		})
		
	case "POST":
		var payload struct {
			Seconds int64 `json:"seconds"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.writeError(w, http.StatusBadRequest, "invalid JSON payload")
			return
		}
		
		if s.db.SetTTL(key, payload.Seconds) {
			s.writeJSON(w, http.StatusOK, map[string]string{"message": "TTL set successfully"})
		} else {
			s.writeError(w, http.StatusNotFound, "key not found")
		}
	}
}

func (s *HTTPServer) handleExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	
	exists := s.db.Exists(key)
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"key":    key,
		"exists": exists,
	})
}

func (s *HTTPServer) handleBulkGet(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Keys []string `json:"keys"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	
	response := s.stringCommands.MGet(payload.Keys)
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"keys":   payload.Keys,
		"values": response.Data,
	})
}

func (s *HTTPServer) handleBulkSet(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Data map[string]string `json:"data"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	
	response := s.stringCommands.MSet(payload.Data)
	if response.Success {
		s.writeJSON(w, http.StatusOK, map[string]interface{}{
			"message": "bulk set successful",
			"count":   len(payload.Data),
		})
	} else {
		s.writeError(w, http.StatusInternalServerError, response.Error)
	}
}

func (s *HTTPServer) handleFlushAll(w http.ResponseWriter, r *http.Request) {
	s.db.FlushAll()
	s.writeJSON(w, http.StatusOK, map[string]string{"message": "database flushed"})
}

// Utility functions
func (s *HTTPServer) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *HTTPServer) writeError(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, map[string]string{"error": message})
}
