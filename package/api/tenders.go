package api

import (
	"avitoproject/package/db"
	"avitoproject/package/entitydescripts"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// получение списка тендеров определенного типа (назначения)
func GetTenders(rw http.ResponseWriter, req *http.Request) {

	serviceType := req.URL.Query().Get("serviceType")

	query := "SELECT tender_id, name, description, service_type, status, organization_id, author_id, version FROM tender"
	var args []interface{} //некоторое колво переменных неявного типа, как auto
	if serviceType != "" {
		query += " WHERE service_type = $1"
		args = append(args, serviceType)
	}

	conn, err := db.Connect() //NOTE засунуть все эти проверялки в одну функцию
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(req.Context(), query, args...)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tenders []entitydescripts.Tender //слайс тендеров (сколько-то)
	for rows.Next() {
		var tender entitydescripts.Tender
		if err := rows.Scan(&tender.TenderID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.AuthorID, &tender.Version); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		tenders = append(tenders, tender)
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(tenders); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// создание нового тендера
func CreateTender(rw http.ResponseWriter, req *http.Request) {

	var newTender entitydescripts.Tender

	if err := json.NewDecoder(req.Body).Decode(&newTender); err != nil {
		http.Error(rw, "Invalid input", http.StatusBadRequest)
		return
	}

	newTender.TenderID = uuid.New()

	conn, err := db.Connect()
	if err != nil {
		http.Error(rw, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	query := `INSERT INTO tender (tender_id, name, service_type, description, status, organization_id, author_id, version) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING tender_id`
	err = conn.QueryRow(context.Background(), query, newTender.TenderID, newTender.Name, newTender.Description, newTender.ServiceType, newTender.Status, newTender.OrganizationID, newTender.AuthorID, newTender.Version).Scan(&newTender.TenderID)
	if err != nil {
		http.Error(rw, "Failed to create tender:"+err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(newTender)
}

// вывод тендеров конкретного автора
func GetSmbdTenders(rw http.ResponseWriter, req *http.Request) {

	username := req.URL.Query().Get("author_id")

	if username == "" {
		http.Error(rw, "Username required", http.StatusBadRequest)
		return
	}

	conn, err := db.Connect()
	if err != nil {
		http.Error(rw, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	query := `SELECT tender_id, name, description, service_type, status, organization_id, author_id, version
              FROM tender WHERE author_id = $1`
	rows, err := conn.Query(context.Background(), query, username)
	if err != nil {
		http.Error(rw, "Failed to execute query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tenders []entitydescripts.Tender
	for rows.Next() {
		var tender entitydescripts.Tender
		if err := rows.Scan(&tender.TenderID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.AuthorID, &tender.Version); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		tenders = append(tenders, tender)
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(tenders)
}

// редактирование тендера
func EditTender(rw http.ResponseWriter, req *http.Request) {
	tenderID := req.URL.Path[len("/api/tenders/") : len(req.URL.Path)-len("/edit")]
	var updateData struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	if err := json.NewDecoder(req.Body).Decode(&updateData); err != nil {
		http.Error(rw, "Invalid input", http.StatusBadRequest)
		return
	}

	conn, err := db.Connect()
	if err != nil {
		http.Error(rw, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	query := `UPDATE tender SET name = COALESCE($1, name), description = COALESCE($2, description) WHERE tender_id = $3 RETURNING tender_id, name, description, service_type, status, organization_id, author_id, version`
	var tender entitydescripts.Tender
	err = conn.QueryRow(context.Background(), query, updateData.Name, updateData.Description, tenderID).Scan(&tender.TenderID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.AuthorID)
	if err != nil {
		http.Error(rw, "Failed to update tender", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(tender)
}
