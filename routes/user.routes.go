package routes

import (
	"github.com/LucioSchiavoni/tas-api/controllers"
	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router) {
	router.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/userById", controllers.GetUserById).Methods("GET")
	router.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	router.HandleFunc("/auth", controllers.ProtectedHandler).Methods("GET")
}
