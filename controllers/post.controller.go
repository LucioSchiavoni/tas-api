package controllers

import (
	"net/http"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion crear post"))
}
