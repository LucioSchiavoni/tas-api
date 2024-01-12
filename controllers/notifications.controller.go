package controllers

import "net/http"

// Notification
func CreateNotification(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para crear la notificacion"))
}

func GetNotificationByUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para traer las notificaciones por el id del usuario que las recibe (UserID)"))
}

// Likes
func CreateLike(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para crear el like"))
}

func GetLikesByIdPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para traer todos los likes por el id del post"))
}

// Comments
func CreateComments(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para crear los comentarios"))
}

func GetCommentsByIdPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para traer los comentarios por el id del post"))
}
