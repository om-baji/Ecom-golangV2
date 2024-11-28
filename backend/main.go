package main

import (
	"backend/routes"
	"log"
	"net/http"
)

func main() {

	log.Fatal(http.ListenAndServe(":4000", routes.UserRouter()))
}
