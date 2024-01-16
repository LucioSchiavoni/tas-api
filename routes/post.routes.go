package routes

import (
	"github.com/LucioSchiavoni/tas-api/controllers"
	"github.com/gorilla/mux"
)

func PostRoutes(router *mux.Router) {

	router.HandleFunc("/post", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/post/{user_id}", controllers.GetPostByIdUser).Methods("GET")
	router.HandleFunc("/AllPost", controllers.GetAllPost).Methods("GET")
	router.HandleFunc("/post/{id}", controllers.DeletePost).Methods("DELETE")
}
