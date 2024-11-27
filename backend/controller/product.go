package controller

import (
	"backend/db"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucsky/cuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection
var cartCollection *mongo.Collection

func init() {
	productCollection = db.Database.Collection("products")
	cartCollection = db.Database.Collection("carts")
}

func addProduct(product models.Product) error {
	product.Id = cuid.New()

	_, err := productCollection.InsertOne(context.Background(), product)
	if err != nil {
		log.Printf("Error adding product: %v\n", err)
		return err
	}

	log.Println("Product successfully added")
	return nil
}

func productToCart(cartId string, productId string, quantity int) error {
	filter := bson.M{"_id": cartId}

	var product models.Product
	if err := productCollection.FindOne(context.Background(), bson.M{"_id": productId}).Decode(&product); err != nil {
		log.Printf("Product with ID %s not found: %v\n", productId, err)
		return fmt.Errorf("product not found")
	}

	cartItem := models.CartItem{
		ProductId: productId,
		Quantity:  quantity,
	}

	update := bson.M{"$push": bson.M{"products": cartItem}}

	_, err := cartCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating cart with ID %s: %v\n", cartId, err)
		return fmt.Errorf("could not update cart")
	}

	log.Printf("Product %s added to cart %s\n", productId, cartId)
	return nil
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("Invalid product data: %v\n", err)
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	if err := addProduct(product); err != nil {
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product added successfully!"})
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := mux.Vars(r)

	cartId := payload["cartId"]

	type CartRequest struct {
		ProductId string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}

	var cartRequest CartRequest

	err := json.NewDecoder(r.Body).Decode(&cartRequest)

	if err != nil {
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}

	err = productToCart(cartId, cartRequest.ProductId, cartRequest.Quantity)

	if err != nil {
		http.Error(w, "Failed to add product", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Product Added to Cart"})

}
