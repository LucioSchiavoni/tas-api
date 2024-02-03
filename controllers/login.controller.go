package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/middlewares"
	"github.com/LucioSchiavoni/tas-api/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("The request body is %v\n", r.Body)

	var loginCredentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginCredentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request body")
		return
	}

	var user models.User
	result := db.DB.Where("email = ?", loginCredentials.Email).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	if CheckPasswordHash(user.Password, loginCredentials.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	tokenString, err := middlewares.CreateToken(user.ID, user.Username, user.Email, user.Image, user.ImageBg, user.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error generando el token")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&tokenString)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Error de Authorization")
		return
	}
	tokenString = tokenString[len("Bearer "):]

	claims, err := middlewares.VerifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	username, ok := claims["username"].(string)
	email, ok := claims["email"].(string)
	image, ok := claims["image"].(string)
	imageBg, ok := claims["image_bg"].(string)
	description, ok := claims["description"].(string)
	userID, ok := claims["id"].(float64)

	responseData := map[string]interface{}{
		"id":          userID,
		"username":    username,
		"email":       email,
		"image":       image,
		"image_bg":    imageBg,
		"description": description}

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error al obtener el username del token")
		return
	}

	json.NewEncoder(w).Encode(&responseData)

}
