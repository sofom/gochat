package utils

/*
	Rest API for chat V1
+	POST /users -- add new user
+	POST /chats -- create a chat for users

+	POST /chats/messages -- send a message

+	GET /chats/<chat_id> -- []Messages from chat history
	DELETE /chats/<chat_id> -- remove the chat

+	GET /chats/<chat_id>/users -- list of the chat members
	POST /chats/<chat_id>/users -- update the chat memebrs list
*/

import (
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServerError struct {
	StatusCode int
	Payload    string
}

func (e *ServerError) Error() string {
	return e.Payload
}

func checkUsers(userIDs []primitive.ObjectID) ([]User, error) {
	var usersList []User

	for _, userID := range userIDs {
		user, senderErr := GetUser(userID)
		if senderErr != nil {
			err := ServerError{StatusCode: http.StatusNotFound}
			err.Payload = fmt.Sprintf("User with UserID %v not found", userID)
			return nil, &err
		}
		usersList = append(usersList, user)
	}

	return usersList, nil
}

func CreateChat(chat *Chat) error {
	_, err := checkUsers(chat.Users)
	if err != nil {
		return err
	}

	_, chatErr := GetChatByUsers(chat.Users)
	if chatErr == nil {
		return errors.New("Chat is already exists")
	}

	if chat.Name == "" {
		chat.Name = "private"
	}
	chat.Create()
	return nil
}

func GetChatMessages(ChatID string) ([]Message, error) {
	chat, err := GetChat(ChatID)
	if err != nil {
		return nil, err
	}

	messages, err := chat.Messages()
	return messages, err
}

func AddChatMessage(newMessage Message) error {
	users := []primitive.ObjectID{newMessage.SenderID, newMessage.ReceiverID}
	_, err := checkUsers(users)
	if err != nil {
		return err
	}

	chat, chatErr := GetChatByUsers(users)
	if chatErr != nil {
		return chatErr
	}
	newMessage.ChatID = chat.ID
	err = newMessage.Create()
	return err
}

func GetChatList(users []primitive.ObjectID, chatType string) ([]Chat, error) {
	_, err := checkUsers(users)
	if err != nil {
		return nil, err
	}
	return getChatList(users, chatType)
}
