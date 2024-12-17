package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func ChatbotHandler(w http.ResponseWriter, r *http.Request) {

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse incoming request
	var req Message
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Prepare payload for Hugging Face
	payload := Payload{
		Model: "microsoft/Phi-3.5-mini-instruct",
		Messages: []Message{
			{
				Role:    "user",
				Content: req.Content,
			},
		},
		MaxToken: 500,
		Stream:   false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
		return
	}

	// Create request to Hugging Face
	huggingReq, err := http.NewRequest("POST", HUGGING_URL+"microsoft/Phi-3.5-mini-instruct/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	huggingReq.Header.Set("Authorization", "Bearer "+TOKEN)
	huggingReq.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(huggingReq)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response
	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	// Send response back to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Choices[0].Message)
}
