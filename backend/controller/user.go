package controller

import (
	"backend/db"
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
	"github.com/lucsky/cuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection
var accountCollection *mongo.Collection
var client *mongo.Client

func checkNilError(err error) {
	if err != nil {
		fmt.Printf("Something went wrong! %w", err)
	}
}

func init() {
	userCollection = db.Database.Collection("users")
	accountCollection = db.Database.Collection("accounts")
}

func addUser(user models.User) error {

	password, success := helper.HashPassword(user.Password)

	if !success {
		log.Fatal("Something went wrong!")
		return errors.New("Hashing")
	}

	user.Password = password
	user.IsAdmin = false
	user.CartId = cuid.New()
	user.AccountNumber = cuid.New()
	account := &models.Account{
		AccountNumber: user.AccountNumber,
		Balance:       1000,
	}

	accountReponse, err := accountCollection.InsertOne(context.Background(), account)

	checkNilError(err)

	response, err := userCollection.InsertOne(context.TODO(), user)

	checkNilError(err)

	fmt.Println("Insertion succesfull! and Account Created", response.InsertedID, accountReponse.InsertedID)

	return nil

}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var credentials models.Credentials
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&credentials)
	checkNilError(err)

	filter := bson.M{"email": credentials.Email}

	response := userCollection.FindOne(context.Background(), filter, nil)

	response.Decode(&user)

	// fmt.Println(response.Raw())

	success := helper.VerifyPassword(user.Password, credentials.Password)

	if !success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bson.M{"msg": "Login Failed (wrong password)"})
		return
	}

	expirationTime := time.Now().Add(time.Minute * 15)

	claims := &models.Claims{
		Username: credentials.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	result := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	token, err := result.SignedString([]byte(string(secret)))

	checkNilError(err)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})
	json.NewEncoder(w).Encode(bson.M{"msg": "Login Successfull"})
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	err := addUser(user)
	if err != nil {
		log.Fatal("Somthing went wrong!")
		return
	}

	expirationTime := time.Now().Add(time.Minute * 15)

	claims := &models.Claims{
		Username: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		http.Error(w, "JWT_SECRET is not set", http.StatusInternalServerError)
		return
	}

	result := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := result.SignedString([]byte(secret))

	checkNilError(err)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

}
