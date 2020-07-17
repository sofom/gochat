package utils

/*
	Rest API for chat V1
	GET /chats/<user_id> - []Messages for the user with the user_id
	POST /chats - Send a message
		payload
		{
			text    string
			from       string
			to         string
		}

	GET /groups/<group_id> -- []Messages from chat history
	POST /groups/<group_id> -- Send new message
		payload
		{
			Message      string
			GroupID      string
		}

	POST /groups/<group_id>/users/<user_id> -- add the user to the group chat
	GET /groups/<group_id>/users/<user_id> -- list of the group members
	DELETE /groups/<group_id>/users/<user_id> -- remove the user from the chat memebrs list
*/

import (
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServerError struct {
	StatusCode int
	payload    string
}

func (e *ServerError) Error() string {
	return e.payload
}

func updateChat(message Message) {
	// TODO: somehow reciver chat should be updated to show new msg
}

func checkUsers(senderID, reciverID primitive.ObjectID) error {
	err := ServerError{StatusCode: http.StatusNotFound}
	if senderErr := GetUser(senderID); senderErr != nil {
		err.payload = fmt.Sprintf("User with UserID %v not found", senderID)
		return &err
	}

	if reciverErr := GetUser(reciverID); reciverErr != nil {
		err.payload = fmt.Sprintf("User with UserID %v not found", reciverID)
		return &err
	}

	return nil
}

func GetChatMessages(UserID string) ([]Message, error) {
	list := make([]Message, 2)
	err := ServerError{StatusCode: http.StatusNotFound, payload: "Chat history is not found"}
	return list, &err
}

//AddChatMessage To many dock strings
func AddChatMessage(newMessage Message) error {

	if err := checkUsers(newMessage.SenderID, newMessage.ReciverID); err != nil {
		return err
	}

	newMessage.Timestamp = time.Now()
	err := newMessage.Create()
	updateChat(newMessage)
	return err
}
