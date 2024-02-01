package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
)

func UploadFile(w http.ResponseWriter, r *http.Request, fieldName string) (string, error) {
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		w.Write([]byte("Error al subir la imagen"))
		return "", err
	}
	defer file.Close()

	fileExtension := filepath.Ext(fileHeader.Filename)

	randomName := fmt.Sprintf("upload-%d%s", time.Now().UnixNano(), fileExtension)
	credPath, err := filepath.Abs("socialapp-go-39fc0a7f2ed2.json")
	if err != nil {
		fmt.Printf("Error del json: %s", err.Error())
		return "", err
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credPath))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error al crear el cliente de Google Cloud Storage: %s", err.Error())
		return "", err
	}
	defer client.Close()

	bucketName := "social-go"
	objectName := "images/" + randomName

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	wc := object.NewWriter(ctx)

	if _, err = io.Copy(wc, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error al copiar el archivo al bucket de GCS: %s", err.Error())
		return "", err
	}

	if err := wc.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error al cerrar el escritor del archivo en GCS: %s", err.Error())
		return "", err
	}

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	return imageURL, nil
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
	w.Header().Set("Content-Type", "json/application")
	var users []models.User
	db.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)

}

//Borrar usuario

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User
	db.DB.First(&user, params["id"])
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	db.DB.Delete(&user)
	w.WriteHeader(http.StatusOK)
}

// Obtener usuario por su id
func GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")
	params := mux.Vars(r)
	var user models.User
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	db.DB.Model(&user).Association("Posts").Find(&user.Post)
	json.NewEncoder(w).Encode(&user)
}
