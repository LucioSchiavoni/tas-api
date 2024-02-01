package routes

import (
	"github.com/LucioSchiavoni/tas-api/controllers"
	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router) {

	router.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	router.HandleFunc("/auth", controllers.ProtectedHandler).Methods("GET")
	router.HandleFunc("/allUser", controllers.GetAllUser).Methods("GET")
	router.HandleFunc("/user/{id}", controllers.DeleteUserHandler).Methods("DELETE")
	router.HandleFunc("/user/{id}", controllers.GetUserById).Methods("GET")
	router.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
}
