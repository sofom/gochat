package v1

import (
	"chat/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func respondWithJSON(w http.ResponseWriter, code int, Payload interface{}) {
	response, _ := json.Marshal(Payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, err error) {
	serverError, ok := err.(*utils.ServerError)
	statusCode := http.StatusBadRequest
	if ok {
		statusCode = serverError.StatusCode
	}
	respondWithJSON(w, statusCode, err.Error())
	log.Print(err.Error())
}

func AddNewUser(w http.ResponseWriter, r *http.Request) {
	var user utils.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, errors.New("Invalid request Payload"))
		return
	}
	defer r.Body.Close()

	err := user.Create()
	if err != nil {
		respondWithError(w, err)
		log.Print(err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func CreateChat(w http.ResponseWriter, r *http.Request) {
	var chat utils.Chat

	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		respondWithError(w, errors.New("Invalid request Payload"))
		return
	}
	defer r.Body.Close()

	err := utils.CreateChat(&chat)
	if err != nil {
		respondWithError(w, err)
		log.Print(err)
		return
	}
	respondWithJSON(w, http.StatusCreated, chat)
}

func ChatList(w http.ResponseWriter, r *http.Request) {
	var users []primitive.ObjectID
	userIDQueryParam := r.FormValue("user")
	if userIDQueryParam != "" {
		userID, err := primitive.ObjectIDFromHex(userIDQueryParam)
		if err != nil {
			respondWithError(w, err)
			return
		}
		users = append(users, userID)
	}

	chatType := r.FormValue("type")
	if chatType != "" && chatType != "private" && chatType != "group" {
		respondWithError(w, errors.New("Invalid query param `type` value"))
		return
	}

	chats, err := utils.GetChatList(users, chatType)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, chats)
}

func ChatMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messages, err := utils.GetChatMessages(vars["id"])
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, messages)
}

func DeleteChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := utils.DeleteChat(vars["id"])
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Deleted")
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var message utils.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		respondWithError(w, errors.New("Invalid request Payload"))
		return
	}
	defer r.Body.Close()

	err := utils.AddChatMessage(message)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, message)
}

func ChatUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chat, err := utils.GetChat(vars["id"])
	if err != nil {
		respondWithError(w, err)
		return
	}

	var users []utils.User
	for _, UserID := range chat.Users {
		user, _ := utils.GetUser(UserID)
		users = append(users, user)
	}

	respondWithJSON(w, http.StatusOK, users)
}

func UpdateChatUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chat, err := utils.GetChat(vars["id"])
	if err != nil {
		respondWithError(w, err)
		return
	}

	type payload struct {
		IDs []primitive.ObjectID `json:"users"`
	}

	var users payload

	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		respondWithError(w, err)
		return
	}
	defer r.Body.Close()

	if len(users.IDs) < 2 {
		respondWithError(w, errors.New("Invalid request Payload. Min number of chat users is 2"))
		return
	}

	err = chat.AddUsers(users.IDs)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, chat)
}

func Handlers(router *mux.Router) {
	router.HandleFunc("/users", AddNewUser).Methods("POST")

	router.HandleFunc("/chats", CreateChat).Methods("POST")

	router.HandleFunc("/chats", ChatList).Methods("GET")
	router.HandleFunc("/chats/{id:[0-9a-z]+}", DeleteChat).Methods("DELETE")

	router.HandleFunc("/chats/{id}/messages", ChatMessages).Methods("GET")
	router.HandleFunc("/chats/messages", SendMessage).Methods("POST")

	router.HandleFunc("/chats/{id}/users", ChatUsers).Methods("GET")
	router.HandleFunc("/chats/{id}/users", UpdateChatUsers).Methods("POST")
}
