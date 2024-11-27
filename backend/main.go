package main

import (
	"backend/routes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))

	log.Fatal(http.ListenAndServe(":4000", routes.UserRouter()))
}
