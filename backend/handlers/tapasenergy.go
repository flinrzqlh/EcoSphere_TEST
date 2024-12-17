package handlers

import (
    "bytes"
    "encoding/csv"
    "encoding/json"
    "net/http"
    "strconv"
    // "strings"
)

func EnergyAnalysisHandler(w http.ResponseWriter, r *http.Request) {

    // Handle preflight requests
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse multipart form
    err := r.ParseMultipartForm(10 << 20) // Max 10 MB
    if err != nil {
        http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
        return
    }

    // Get the file from the request
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Read the CSV file
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        http.Error(w, "Could not read CSV file", http.StatusInternalServerError)
        return
    }

    // Aggregate energy consumption by appliance
    applianceConsumption := make(map[string]float64)

    // Skip header and parse data
    var table Table
    for i, record := range records {
        if i == 0 {
            continue // skip header
        }
        table.Date = append(table.Date, record[0])
        table.Time = append(table.Time, record[1])
        table.Appliance = append(table.Appliance, record[2])
        table.EnergyConsumption = append(table.EnergyConsumption, record[3])
        table.Room = append(table.Room, record[4])
        table.Status = append(table.Status, record[5])

        // Aggregate consumption
        consumption, err := strconv.ParseFloat(record[3], 64)
        if err == nil {
            applianceConsumption[record[2]] += consumption
        }
    }

    // Find least and most consuming appliances
    var leastConsumingAppliance string
    var mostConsumingAppliance string
    var leastConsumption float64 = float64(^uint(0) >> 1)  // Max possible int value
    var maxConsumption float64 = 0

    for appliance, consumption := range applianceConsumption {
        if consumption < leastConsumption {
            leastConsumption = consumption
            leastConsumingAppliance = appliance
        }

        if consumption > maxConsumption {
            maxConsumption = consumption
            mostConsumingAppliance = appliance
        }
    }

    // Prepare payload for Hugging Face
    payload := EnergyAnalysisPayload{
        Query: "Identify the most and least energy consuming appliances from this dataset.",
        Table: table,
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
        return
    }

    // Create request to Hugging Face
    huggingReq, err := http.NewRequest("POST", HUGGING_URL+"google/tapas-base-finetuned-wtq", bytes.NewBuffer(payloadBytes))
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

    // Prepare response
    response := EnergyAnalysisResponse{
        LeastConsumingAppliance: leastConsumingAppliance,
        MostConsumingAppliance:  mostConsumingAppliance,
        LeastConsumptionValue:   leastConsumption,
        MostConsumptionValue:    maxConsumption,
    }

    // Send JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}