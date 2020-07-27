package main

import (
	"chat/utils"
	v1 "chat/v1"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func app(w http.ResponseWriter, r *http.Request) {
	log.Print("health")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func main() {
	client := utils.SetupDBConnection()
	log.Print("Start server")
	router := mux.NewRouter()

	v1.Handlers(router.PathPrefix("/api/v1").Subrouter())

	router.HandleFunc("/api/health", app)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
	defer client.Disconnect(context.TODO())
}
