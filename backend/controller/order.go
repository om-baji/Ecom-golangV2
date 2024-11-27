package controller

import (
	"backend/db"
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection

func init() {
	orderCollection = db.Database.Collection("orders")
}

func checkOut(cartId string) {

	cartFilter := bson.M{"_id": cartId}

	cart := cartCollection.FindOne(context.Background(), cartFilter, nil)

	var order models.Order

	cart.Decode(&order)

}
