package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id       string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string   `json:"name,omitempty" bson:"name,omitempty"`
	Email    string   `json:"email,omitempty" bson:"email,omitempty"`
	Password string   `json:"password,omitempty" bson:"password,omitempty"`
	Account  *Account `json:"account,omitempty" bson:"account,omitempty"`
	CartId   string   `json:"cartId" bson:"cartId"`
	IsAdmin  bool     `json:isAdmin" bson:"isAdmin"`
}

type Account struct {
	AccountNumber string     `json:"accountNumber" bson:"accountNumber"`
	Balance       int64      `json:"balance" bson:"balance"`
	UserId        string     `json:"userId" bson:"userId"`
	Payments      []*Payment `json:"payments,omitempty" bson:"payments,omitempty"`
}

type Payment struct {
	Id     string    `json:"transactionId" bson:"transactionId"`
	From   string    `json:"from" bson:"from"`
	To     string    `json:"to" bson:"to"`
	Amount int64     `json:"amount" bson:"amount"`
	Time   time.Time `json:"time" bson:"time"`
	Status string    `json:"status" bson:"status"`
}

type Cart struct {
	Id        string     `json:"id" bson:"_id,omitempty"`
	UserId    string     `json:"userId" bson:"userId"`
	Products  []CartItem `json:"products" bson:"products"`
	UpdatedAt time.Time  `json:"updatedAt" bson:"updatedAt"`
}

type CartItem struct {
	ProductId string `json:"productId" bson:"productId"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}

type Order struct {
	Id         string    `json:"orderId" bson:"orderId"`
	UserId     string    `json:"userId" bson:"userId"`
	Products   Cart      `json:"products" bson:"products"`
	TotalPrice int64     `json:"totalPrice" bson:"totalPrice"`
	Status     string    `json:"status" bson:"status"`
	Time       time.Time `json:"time" bson:"time"`
}

type Product struct {
	Id          string `json:"productId" bson:"productId"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Price       int64  `json:"price" bson:"price"`
	Stock       int    `json:"stock" bson:"stock"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
