package controllers

import "net/http"

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para crear la notificacion"))
}
