package api

import (
	"avitoproject/package/db"
	"avitoproject/package/entitydescripts"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func CreateBid(rw http.ResponseWriter, req *http.Request) {
	var newBid entitydescripts.Offer

	if err := json.NewDecoder(req.Body).Decode(&newBid); err != nil {
		http.Error(rw, "Invalid input", http.StatusBadRequest)
		return
	}

	newBid.OfferID = uuid.New()

	conn, err := db.Connect()
	if err != nil {
		http.Error(rw, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	query := `INSERT INTO offer (offer_id, name, description, status, tender_id, organization_id, author_id) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING offer_id`
	err = conn.QueryRow(context.Background(), query, newBid.OfferID, newBid.Name, newBid.Description, newBid.Status, newBid.TenderID, newBid.OrganizationID, newBid.AuthorID).Scan(&newBid.OfferID)
	if err != nil {
		http.Error(rw, "Failed to create bid"+err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(newBid)
}

func GetSmbdBids(rw http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("username")
	if username == "" {
		http.Error(rw, "Username is required", http.StatusBadRequest)
		return
	}

	conn, err := db.Connect()
	if err != nil {
		http.Error(rw, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	query := `SELECT offer_id, name, description, status, tender_id, organization_id, author_id FROM offer WHERE author_id = $1`
	rows, err := conn.Query(context.Background(), query, username)
	if err != nil {
		http.Error(rw, "Failed to retrieve bids", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bids []entitydescripts.Offer
	for rows.Next() {
		var bid entitydescripts.Offer
		if err := rows.Scan(&bid.OfferID, &bid.Name, &bid.Description, &bid.Status, &bid.OrganizationID, &bid.AuthorID); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		bids = append(bids, bid)
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(bids)
}

func GetBidsForTender(rw http.ResponseWriter, req *http.Request) {
	tenderId := req.URL.Query().Get("tenderId")
	if tenderId == "" {
		http.Error(rw, "Tender ID is required", http.StatusBadRequest)
		return
	}

	conn, err := db.Connect()
	if err != nil {
		http.Error(rw, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	query := `SELECT offer_id, name, description, status, tender_id, organization_id, author_id FROM offer WHERE tender_id = $1`
	rows, err := conn.Query(context.Background(), query, tenderId)
	if err != nil {
		http.Error(rw, "Failed to retrieve bids", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bids []entitydescripts.Offer
	for rows.Next() {
		var bid entitydescripts.Offer
		if err := rows.Scan(&bid.OfferID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.OrganizationID, &bid.AuthorID); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		bids = append(bids, bid)
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(bids)
}

func EditBid(rw http.ResponseWriter, req *http.Request) {
	bidId := req.URL.Query().Get("offerId")
	if bidId == "" {
		http.Error(rw, "Bid ID is required", http.StatusBadRequest)
		return
	}

	var updatedBid entitydescripts.Offer
	if err := json.NewDecoder(req.Body).Decode(&updatedBid); err != nil {
		http.Error(rw, "Invalid input", http.StatusBadRequest)
		return
	}

	conn, err := db.Connect()
	if err != nil {
		http.Error(rw, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	query := `UPDATE offer SET name = $1, description = $2 WHERE offer_id = $3`
	_, err = conn.Exec(context.Background(), query, updatedBid.Name, updatedBid.Description, bidId)
	if err != nil {
		http.Error(rw, "Failed to update bid", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(updatedBid)
}
