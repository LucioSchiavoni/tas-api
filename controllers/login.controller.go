package controllers

import "net/http"

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion del login"))
}
