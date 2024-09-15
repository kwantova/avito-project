package api

import (
	"avitoproject/package/db"
	"avitoproject/package/entitydescripts"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type ChangeStatusT struct {
	CurrentStatus string `json:"status"`
}

func UpdateTStatus(w http.ResponseWriter, req *http.Request) {
	tIDStr := req.URL.Query().Get("TenderID")
	if tIDStr == "" {
		http.Error(w, "TenderID is required", http.StatusBadRequest)
		return
	}
	tID, err := uuid.Parse(tIDStr)
	if err != nil {
		http.Error(w, "Invalid TenderID format", http.StatusBadRequest)
		return
	}

	var request ChangeStatusT

	if json.NewDecoder(req.Body).Decode(&request) != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	AllPossibleStats := map[string]bool{
		"CREATED":   true,
		"PUBLISHED": true,
		"CLOSED":    true}

	if !AllPossibleStats[request.CurrentStatus] {
		http.Error(w, "Invalid request: ", http.StatusBadRequest)
		return
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close(context.Background())

	query := "UPDATE tender SET status = $1 WHERE tender_id = $2 RETURNING status, tender_id"
	var updatedTender entitydescripts.Tender
	err = db.QueryRow(context.Background(), query, request.CurrentStatus, tID).Scan(&updatedTender.Status, &updatedTender.TenderID)
	if err != nil {
		http.Error(w, "Failed to update tender status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTender); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
