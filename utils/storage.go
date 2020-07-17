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
	ReciverID primitive.ObjectID `bson:"reciver_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

type User struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type Group struct {
	ID      primitive.ObjectID `bson:"_id"`
	Name    string             `bson:"name"`
	Members []User             `bson:"members"`
}

var client *mongo.Client

func GetUser(UserID primitive.ObjectID) error {
	filter := bson.M{"_id": UserID}

	var result User
	users := client.Database("chat").Collection("users")
	err := users.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return err
}


func (m Message) Create() error {
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
