package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"concerts/db"
	"concerts/handlers"
)

func main() {
    if _, err := db.Init(); err != nil {
        log.Fatalf("db init failed: %v", err)
    }

    r := mux.NewRouter()
    r.Use(corsMiddleware)

    // Handle CORS preflight for all routes
    r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusNoContent)
    }).Methods(http.MethodOptions)

    // Health
    r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)

    // Auth routes
    r.HandleFunc("/register", handlers.Register).Methods(http.MethodPost)
    r.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)

    // Concerts (protected)
    concerts := r.PathPrefix("/concerts").Subrouter()
    concerts.Use(handlers.RequireAuth)
    concerts.HandleFunc("", handlers.ListConcerts).Methods(http.MethodGet)
    concerts.HandleFunc("/", handlers.ListConcerts).Methods(http.MethodGet)
    concerts.HandleFunc("", handlers.CreateConcert).Methods(http.MethodPost)
    concerts.HandleFunc("/", handlers.CreateConcert).Methods(http.MethodPost)
    concerts.HandleFunc("/{id}", handlers.DeleteConcert).Methods(http.MethodDelete)

    srv := &http.Server{
        Addr:              getAddr(),
        Handler:           r,
        ReadTimeout:       10 * time.Second,
        ReadHeaderTimeout: 10 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }
    log.Printf("server listening on %s", srv.Addr)
    log.Fatal(srv.ListenAndServe())
}

func getAddr() string {
    if p := os.Getenv("PORT"); p != "" {
        return ":" + p
    }
    return ":8080"
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}


