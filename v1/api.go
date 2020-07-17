package v1

import (
	"chat/utils"

	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ChatList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messages, err := utils.GetChatMessages(vars["id"])
	if err != nil {
		serverError := err.(*utils.ServerError)
		respondWithJSON(w, serverError.StatusCode, serverError.Error())
		log.Print(err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, messages)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var message utils.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := utils.AddChatMessage(message)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, err)
		log.Panic(err)
		return
	}
	respondWithJSON(w, http.StatusOK, message)
}

func Handlers(router *mux.Router) {
	router.HandleFunc("/chats/{id:[0-9]+}", ChatList).Methods("GET")
	router.HandleFunc("/chats", SendMessage).Methods("POST")

	// router.HandleFunc("/products", ProductsHandler)
	// router.HandleFunc("/articles", ArticlesHandler)
}
