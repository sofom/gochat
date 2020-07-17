package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Text      string             `bson:"text"`
	SenderID  primitive.ObjectID `bson:"sender_id"`
	receiveID primitive.ObjectID `bson:"receive_id"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp,omitempty"`
	ChatID    primitive.ObjectID `bson:"chat_id"`
}

type User struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
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
	if err != nil {
		log.Print(err)
	}

	return user, err
}

func GetChat(users []primitive.ObjectID) (chat Chat, err error) {
	filter := bson.M{"members": bson.D{{"$all", users}}}

	collection := client.Database("chat").Collection("chats")
	err = collection.FindOne(context.TODO(), filter).Decode(&chat)
	if err != nil {
		log.Print(err)
		return chat, err
	}
	return chat, err
}

func (c *Chat) Create() error {
	c.ID = primitive.NewObjectID()
	collection := client.Database("chat").Collection("chats")
	inserted, err := collection.InsertOne(context.TODO(), c)
	log.Printf("New chat %v", inserted)
	return err
}

func DeleteAllChats() {
	collection := client.Database("chat").Collection("chats")
	_, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
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
