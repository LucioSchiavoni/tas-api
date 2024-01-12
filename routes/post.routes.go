package routes

import (
	"github.com/LucioSchiavoni/tas-api/controllers"
	"github.com/gorilla/mux"
)

func PostRoutes(router *mux.Router) {

	router.HandleFunc("/post", controllers.CreatePost).Methods("POST")
}
