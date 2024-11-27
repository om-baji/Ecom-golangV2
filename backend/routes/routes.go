package routes

import (
	"backend/controller"

	"github.com/gorilla/mux"
)

func UserRouter() *mux.Router {

	userRouter := mux.NewRouter()

	userRouter.HandleFunc("/login", controller.Login).Methods("POST")
	userRouter.HandleFunc("/signup", controller.Register).Methods("POST")

	return userRouter

}
