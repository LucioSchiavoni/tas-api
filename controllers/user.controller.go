package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
	"golang.org/x/crypto/bcrypt"
)

func UploadFile(w http.ResponseWriter, r *http.Request, fieldName string) (string, error) {
	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile(fieldName)
	if err != nil {
		w.Write([]byte("Error al subir la imagen"))
		return "", err

	}

	defer file.Close()

	tempFile, err := os.CreateTemp("images", "upload-*.png")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear el archivo temporal"})
		return "", err
	}

	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)

	if err != nil {
		fmt.Println(err)
	}

	return tempFile.Name(), nil

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error de formato FormData"})
		return
	}

	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	user.Image = r.FormValue("image")
	user.ImageBg = r.FormValue("image_bg")
	user.Description = r.FormValue("description")

	// Hash password
	secret := os.Getenv("HASH_PWD")
	password := r.FormValue("password")

	hash, err := HashPassword(secret + password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al generar el hash de la contraseña"})
		return
	}

	user.Password = string(hash)

	// Upload File
	file, _, err := r.FormFile("image")
	fileBg, _, err := r.FormFile("image_bg")

	imagePath, err := UploadFile(w, r, "image")
	imagePathBg, err := UploadFile(w, r, "image_bg")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al subir la imagen"})
	}
	defer file.Close()
	defer fileBg.Close()

	user.Image = imagePath
	user.ImageBg = imagePathBg

	createUser := db.DB.Create(&user)
	err = createUser.Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"Error": err.Error()})
		return
	}

	user.Password = ""

	json.NewEncoder(w).Encode(&user)

}

// Todos los usuarios
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para traer todos los usuarios"))
}

// Obtener usuario por su id
func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Obtener usuario by ID"))
}
