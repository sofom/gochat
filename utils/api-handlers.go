package utils

/*
	Rest API for chat V1
	POST /users -- add new user
	POST /chats -- create a chat for users

	POST /chats/messages -- send a message

	GET /chats/<chat_id> -- []Messages from chat history
	DELETE /chats/<chat_id> -- remove the chat

	GET /chats/<chat_id>/users -- list of the chat members

	POST /chats/<chat_id>/users/<user_id> -- add the user to the chat memebrs list
	DELETE /chats/<chat_id>/users/<user_id> -- remove the user from the chat memebrs list
*/

import (
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServerError struct {
	StatusCode int
	payload    string
}

func (e *ServerError) Error() string {
	return e.payload
}

func checkUsers(userIDs []primitive.ObjectID) ([]User, error) {
	var usersList []User

	for _, userID := range userIDs {
		user, senderErr := GetUser(userID)
		if senderErr != nil {
			err := ServerError{StatusCode: http.StatusNotFound}
			err.payload = fmt.Sprintf("User with UserID %v not found", userID)
			return nil, &err
		}
		usersList = append(usersList, user)
	}

	return usersList, nil
}

func CreateChat(chat *Chat) error {
	_, err := checkUsers(chat.Members)
	if err != nil {
		return err
	}

	_, chatErr := GetChat(chat.Members)
	if chatErr == nil {
		return errors.New("Chat is already exists")
	}

	if chat.Name == "" {
		chat.Name = "private"
	}
	chat.Create()
	return nil
}

func GetChatMessages(UserID string) ([]Message, error) {
	list := make([]Message, 2)
	err := ServerError{StatusCode: http.StatusNotFound, payload: "Chat history is not found"}
	return list, &err
}

func AddChatMessage(newMessage Message) error {
	_, err := checkUsers([]primitive.ObjectID{newMessage.SenderID, newMessage.receiveID})
	if err != nil {
		return err
	}

	err = newMessage.Create()
	return err
}
