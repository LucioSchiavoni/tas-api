package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
)

func UploadFileGoogleCloud(w http.ResponseWriter, r *http.Request, fieldName string) (string, error) {
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		w.Write([]byte("Error al subir la imagen"))
		return "", err
	}
	defer file.Close()

	fileExtension := filepath.Ext(fileHeader.Filename)
	randomName := fmt.Sprintf("upload-%d%s", time.Now().UnixNano(), fileExtension)

	// credPath, err := filepath.Abs("socialapp-go-39fc0a7f2ed2.json")
	// if err != nil {
	// 	fmt.Printf("Error del json: %s", err.Error())
	// 	return "", err
	// }

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(fmt.Sprintf(`{
		"type":"service_account",
		"project_id":"%s",
		"private_key_id":"%s",
		"private_key":"%s",
		"client_email":"%s",
		"client_id":"%s",
		"auth_uri":"https://accounts.google.com/o/oauth2/auth",
		"token_uri":"https://accounts.google.com/o/oauth2/token",
		"auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":"%s",
		"universe_domain": "googleapis.com"
	}`, os.Getenv("PROJECT_ID"), os.Getenv("PRIVATE_KEY_ID"), os.Getenv("PRIVATE_KEY"), os.Getenv("CLIENT_EMAIL"), os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_CERT_URL")))))
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

func UploadFile(w http.ResponseWriter, r *http.Request, fieldName string) (string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", fmt.Errorf("error parsing form: %v", err)
	}

	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", fmt.Errorf("error getting form file: %v", err)
	}
	defer file.Close()

	imageDir := "images"
	filePath := filepath.Join(imageDir, fileHeader.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("error copying file: %v", err)
	}

	urlPath := os.Getenv("URL_PATH")
	if urlPath == "" {
		urlPath = "http://localhost:8080"
	}

	fileURL := fmt.Sprintf("%s/images/%s", strings.TrimRight(urlPath, "/"), fileHeader.Filename)

	return fileURL, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error de formato FormData", http.StatusBadRequest)
		return
	}

	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Description = r.FormValue("description")

	// Hash de la contraseña
	password := r.FormValue("password")
	hash, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Error al generar el hash de la contraseña", http.StatusInternalServerError)
		return
	}
	user.Password = string(hash)

	// Subir archivos de imagen
	imagePath, err := UploadFile(w, r, "image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al subir la imagen: %v", err), http.StatusInternalServerError)
		return
	}
	imagePathBg, err := UploadFile(w, r, "image_bg")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al subir la imagen de fondo: %v", err), http.StatusInternalServerError)
		return
	}

	user.Image = imagePath
	user.ImageBg = imagePathBg

	// Verificar si el usuario ya existe
	var existingUser models.User
	result := db.DB.Where("email = ?", user.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		http.Error(w, "Usuario ya registrado", http.StatusConflict)
		return
	}

	// Crear usuario
	err = db.DB.Create(&user).Error
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al crear el usuario: %v", err), http.StatusBadRequest)
		return
	}

	user.Password = "" // No devolver la contraseña en la respuesta

	w.WriteHeader(http.StatusCreated)
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

	db.DB.Unscoped().Delete(&user)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted permanently"))
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

// Editar usuario
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error de formato FormData"})
		return
	}

	err = db.DB.First(&user, params["id"]).Error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no encontrado"})
		return
	}

	// Se cambian los datos solo si hubo un cambio en el campo
	if username := r.FormValue("username"); username != "" {
		user.Username = username
	}

	if password := r.FormValue("password"); password != "" {

		hash, err := HashPassword(password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error al generar el hash de la contraseña"})
			return
		}
		user.Password = string(hash)
	}

	// if image, err := UploadFile(w, r, "image"); err == nil {
	// 	user.Image = image
	// }

	// if imageBg, err := UploadFile(w, r, "image_bg"); err == nil {
	// 	user.ImageBg = imageBg
	// }
	if description := r.FormValue("description"); description != "" {
		user.Description = description
	}

	err = db.DB.Save(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al actualizar el usuario en la base de datos"})
		return
	}

	user.Password = ""
	json.NewEncoder(w).Encode(&user)
}
