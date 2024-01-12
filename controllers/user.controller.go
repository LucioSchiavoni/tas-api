package controllers

import (
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion Create"))
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para traer todos los usuarios"))
}
