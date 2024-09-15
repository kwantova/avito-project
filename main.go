package main

import (
	"avitoproject/package/api"
	"log"
	"net/http"
	"os"
)

func main() {

	addr := os.Getenv("SERVER_ADDRESS")
	log.Printf("Server will start on %s", os.Getenv("SERVER_ADDRESS"))
	log.Printf("Postgres connection string: %s", os.Getenv("POSTGRES_CONN"))
	log.Printf("User: %s", os.Getenv("POSTGRES_USERNAME"))

	/*addr := "0.0.0.0:8080" //костыль
	if addr == "" {
		addr = ":8080"
	}*/

	http.HandleFunc("/api/ping", api.Ping)
	http.HandleFunc("/api/tenders", api.GetTenders)
	http.HandleFunc("/api/tenders/new", api.CreateTender)
	http.HandleFunc("/api/tenders/my", api.GetSmbdTenders)
	http.HandleFunc("/api/tenders/{tenderId}/edit", api.EditTender)
	http.HandleFunc("/api/tenders/status", api.UpdateTStatus)
	http.HandleFunc("/api/bids/status", api.UpdateOStatus)
	http.HandleFunc("/api/bids/new", api.CreateBid)
	http.HandleFunc("/api/bids/my", api.GetSmbdBids)
	http.HandleFunc("/api/bids/{tenderId}/list", api.GetBidsForTender)

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
