package controller

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://localhost:27117"

var userCollection *mongo.Collection
var accountCollection *mongo.Collection
var paymentCollection *mongo.Collection
var cartCollection *mongo.Collection
var orderCollection *mongo.Collection

var client *mongo.Client

func checkNilError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {

	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.Background(), clientOptions)

	checkNilError(err)

	fmt.Println("Database Connected!")

	userCollection = client.Database("ECOM").Collection("user")
	accountCollection = client.Database("ECOM").Collection("account")
	paymentCollection = client.Database("ECOM").Collection("payment")
	cartCollection = client.Database("ECOM").Collection("cart")
	orderCollection = client.Database("ECOM").Collection("order")

}

func CloseConnection() {
	err := client.Disconnect(context.Background())

	checkNilError(err)

	fmt.Println("Db conenction closed")

}
