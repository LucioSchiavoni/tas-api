package routes

import (
	"github.com/LucioSchiavoni/tas-api/controllers"
	"github.com/gorilla/mux"
)

func ChatRouter(router *mux.Router) {
	router.HandleFunc("/message", controllers.CreateMessageByUser).Methods("POST")
	router.HandleFunc("/messsage", controllers.GetMessageByUser).Methods("GET")

}
