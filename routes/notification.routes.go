package routes

import (
	"github.com/LucioSchiavoni/tas-api/controllers"
	"github.com/gorilla/mux"
)

func NotificationRouter(router *mux.Router) {
	// router.HandleFunc("/notification", controllers.CreateNotification).Methods("POST")
	router.HandleFunc("/like", controllers.CreateLike).Methods("POST")
	router.HandleFunc("/{user_id}/like/{id}", controllers.DeleteLike).Methods("DELETE")
	router.HandleFunc("/{post_id}/likes/{user_id}", controllers.GetLikesByIdPost).Methods("GET")
	router.HandleFunc("/comments", controllers.CreateComments).Methods("POST")
	router.HandleFunc("/{post_id}/comments", controllers.GetCommentsByIdPost).Methods("GET")
	router.HandleFunc("/notificationByUser/{user_id}", controllers.GetNotificationByUser).Methods("GET")
}
