package main

import (
    "log"
    "net/http"

    "MembantuKawan/backend/handlers"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    }
}

func main() {
    // Setup routes with CORS middleware
    http.HandleFunc("/chatbot", corsMiddleware(handlers.ChatbotHandler))
    http.HandleFunc("/energy-analysis", corsMiddleware(handlers.EnergyAnalysisHandler))

    // Start server
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}