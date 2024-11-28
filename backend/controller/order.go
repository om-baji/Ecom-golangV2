package controller

import (
	"backend/db"
	"backend/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection

func init() {
	orderCollection = db.Database.Collection("orders")
}

func checkOut(cartId string) error {

	var user models.User
	var cart models.Cart

	userFilter := bson.M{"cartId": cartId}
	cartFilter := bson.M{"_id": cartId}

	userResponse := userCollection.FindOne(context.Background(), userFilter, nil)
	cartResponse := cartCollection.FindOne(context.Background(), cartFilter, nil)

	cartResponse.Decode(&cart)
	userResponse.Decode(&user)

	products := []string{}

	orderEntry := &models.Order{
		Email:    user.Email,
		Products: products,
		Status:   false,
		Time:     time.Now(),
	}

	response, err := orderCollection.InsertOne(context.Background(), orderEntry, nil)

	if err != nil {
		fmt.Println("Something went wrong!")
		return err
	}

	fmt.Println("Order added", response.InsertedID)

	return nil
}
