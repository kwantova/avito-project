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

type ChangeStatusO struct {
	CurrentStatus string `json:"status"`
}

func UpdateOStatus(w http.ResponseWriter, req *http.Request) {
	oIDstr := req.URL.Query().Get("OfferID")
	if oIDstr == "" {
		http.Error(w, "OfferID is required", http.StatusBadRequest)
		return
	}

	oID, err := uuid.Parse(oIDstr)
	if err != nil {
		http.Error(w, "Invalid OfferID format", http.StatusBadRequest)
		return
	}

	var request ChangeStatusO
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	AllPossibleStats := map[string]bool{
		"CREATED":   true,
		"PUBLISHED": true,
		"CANCELLED": true,
	}

	if !AllPossibleStats[request.CurrentStatus] {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	var currentStatus string
	query := "SELECT status FROM offer WHERE offer_id = $1"
	err = db.QueryRow(context.Background(), query, oID).Scan(&currentStatus)
	if err != nil {
		http.Error(w, "Offer not found", http.StatusNotFound)
		return
	}

	validTransition := false

	if currentStatus == "CREATED" && (request.CurrentStatus == "PUBLISHED" || request.CurrentStatus == "CANCELLED") {
		validTransition = true
	} else if currentStatus == "PUBLISHED" && request.CurrentStatus == "CANCELLED" {
		validTransition = true
	}

	if !validTransition {
		http.Error(w, "Invalid status transition", http.StatusBadRequest)
		return
	}

	query = "UPDATE offer SET status = $1 WHERE offer_id = $2 RETURNING status, offer_id"
	var updatedOffer entitydescripts.Offer
	err = db.QueryRow(context.Background(), query, request.CurrentStatus, oID).Scan(&updatedOffer.Status, &updatedOffer.OfferID)
	if err != nil {
		http.Error(w, "Failed to update offer status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOffer)
}
