package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Text       string             `bson:"text"`
	SenderID   primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	ReceiverID primitive.ObjectID `bson:"receiver_id" json:"receiver_id"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp,omitempty"`
	ChatID     primitive.ObjectID `bson:"chat_id"`
}

type User struct {
	ID   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name string             `bson:"name" json:"name"`
}

type Chat struct {
	ID      primitive.ObjectID   `bson:"_id" json:"id,omitempty"`
	Name    string               `bson:"name" json:"name,omitempty"`
	Members []primitive.ObjectID `bson:"members" json:"members"`
}

var client *mongo.Client

func GetUser(UserID primitive.ObjectID) (user User, err error) {
	filter := bson.M{"_id": UserID}

	collection := client.Database("chat").Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

func (u *User) Create() error {
	collection := client.Database("chat").Collection("users")
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"name": u.Name}).Decode(&user)
	if user != (User{}) {
		err := ServerError{StatusCode: http.StatusNotFound}
		err.Payload = fmt.Sprintf("User with Name %v has already registred", u.Name)
		return &err
	}

	u.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(context.TODO(), u)
	return err
}

func GetChat(ChatID string) (chat Chat, err error) {
	id, err := primitive.ObjectIDFromHex(ChatID)
	if err != nil {
		return Chat{}, err
	}

	filter := bson.M{"_id": id}

	collection := client.Database("chat").Collection("chats")
	err = collection.FindOne(context.TODO(), filter).Decode(&chat)
	return chat, err
}

func GetChatByUsers(users []primitive.ObjectID) (chat Chat, err error) {
	filter := bson.M{"members": bson.D{{"$all", users}}}

	collection := client.Database("chat").Collection("chats")
	err = collection.FindOne(context.TODO(), filter).Decode(&chat)
	if err != nil {
		return chat, err
	}
	return chat, err
}

func (c *Chat) Create() error {
	c.ID = primitive.NewObjectID()
	collection := client.Database("chat").Collection("chats")
	_, err := collection.InsertOne(context.TODO(), c)
	return err
}

func (c *Chat) Messages() ([]Message, error) {
	var messageList []Message
	filter := bson.M{"chat_id": c.ID}
	findOptions := options.Find()
	findOptions.SetLimit(20)

	ctx := context.TODO()
	collection := client.Database("chat").Collection("messages")
	cur, err := collection.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var msg Message
		err := cur.Decode(&msg)
		if err != nil {
			err = &ServerError{Payload: err.Error()}
		}

		messageList = append(messageList, msg)
	}
	if err := cur.Err(); err != nil {
		err = &ServerError{Payload: err.Error()}
	}

	return messageList, err
}

func (c *Chat) AddMembers(userIDs []primitive.ObjectID) error {
	_, err := checkUsers(userIDs)
	if err != nil {
		return err
	}

	update := bson.D{
		{"$addToSet", bson.D{
			{"members", bson.D{
				{"$each", userIDs}},
			},
		}},
	}
	collection := client.Database("chat").Collection("chats")
	_, updateErr := collection.UpdateOne(context.TODO(), bson.M{"_id": c.ID}, update)

	if updateErr != nil {
		return updateErr
	}
	c.Members = append(userIDs)
	return nil
}

func DeleteChat(ChatID string) error {
	id, err := primitive.ObjectIDFromHex(ChatID)
	if err != nil {
		return &ServerError{Payload: err.Error()}
	}

	collection := client.Database("chat").Collection("chats")
	_, err = collection.DeleteOne(context.TODO(), bson.M{"chat_id": id})
	if err != nil {
		return &ServerError{Payload: err.Error()}
	}
	return nil
}

func (m Message) Create() error {
	m.Timestamp = time.Now()
	collection := client.Database("chat").Collection("messages")
	_, err := collection.InsertOne(context.TODO(), m)
	return err
}

func SetupDBConnection() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error

	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("DB connection established")
	return client
}
