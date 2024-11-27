package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	once     sync.Once
	err      error
	Database *mongo.Database
)

func init() {
	once.Do(func() {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
		Client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal("Database connection failed:", err)
		}

		Database = Client.Database("ECOM")
		fmt.Println("Database connected!")
	})
}

func CloseConnection() {
	err := Client.Disconnect(context.Background())
	if err != nil {
		fmt.Println("Error closing database connection:", err)
		return
	}
	fmt.Println("Database connection closed")
}
