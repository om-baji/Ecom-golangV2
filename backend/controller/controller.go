package controller

import (
	"backend/helper"
	"backend/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/lucsky/cuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

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

func addUser(user models.User) error {

	password, success := helper.HashPassword(user.Password)

	if !success {
		log.Fatal("Something went wrong!")
		return errors.New("Hashing")
	}

	user.Password = password
	user.IsAdmin = false
	user.Account = &models.Account{
		AccountNumber: cuid.New(),
		Balance:       1000,
		UserId:        user.Id,
	}

	response, err := userCollection.InsertOne(context.TODO(), user)

	checkNilError(err)

	fmt.Println("Insertion succesfull! ", response.InsertedID)

	return nil

}

func Login(w http.ResponseWriter, r *http.Request) (string, error) {
	var credentials models.Credentials
	var user models.User

	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	checkNilError(err)

	filter := bson.M{"username": credentials.Username}

	response := userCollection.FindOne(context.Background(), filter, nil)

	response.Decode(&user)

	success := helper.VerifyPassword(user.Password, credentials.Password)

	if !success {
		w.WriteHeader(http.StatusBadRequest)
		return "", errors.New("wrong password")
	}

	expirationTime := time.Now().Add(time.Minute * 15)

	claims := &models.Claims{
		Username: credentials.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	result := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := result.SignedString(os.Getenv("JWT_SECRET"))

	checkNilError(err)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	return token, nil
}

func Register(w http.ResponseWriter, r *http.Request) (string, error) {

	var user models.User

	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	json.NewDecoder(r.Body).Decode(&user)

	err := addUser(user)
	if err != nil {
		log.Fatal("Somthing went wrong!")
		return "", nil
	}

	expirationTime := time.Now().Add(time.Minute * 15)

	claims := &models.Claims{
		Username: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	result := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := result.SignedString(os.Getenv("JWT_SECRET"))

	checkNilError(err)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	return token, nil

}
