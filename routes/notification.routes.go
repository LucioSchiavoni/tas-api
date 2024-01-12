package routes

import (
	"github.com/LucioSchiavoni/tas-api/controllers"
	"github.com/gorilla/mux"
)

func NotificationRouter(router *mux.Router) {
	router.HandleFunc("/notification", controllers.CreateNotification).Methods("POST")
}
