package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error de formato FormData"})
		return
	}

	description := r.FormValue("description")
	userIDStr := r.FormValue("id")

	file, _, err := r.FormFile("image_post")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if err != nil {
		http.Error(w, "Error al aprsear UserID", http.StatusBadRequest)
		return
	}

	var imagePath string

	if file != nil {
		imagePath, err = UploadFile(w, r, "image_post")

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error al subir la imagen"})
		}
		defer file.Close()
	}

	newPost := models.Post{
		Description: description,
		ImagePost:   imagePath,
		UserID:      uint(userID),
	}

	createPost := db.DB.Create(&newPost)
	err = createPost.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(&post)

}

func GetPostByIdUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")
	params := mux.Vars(r)

	userID := params["user_id"]
	var post []models.Post
	if err := db.DB.Where("user_id = ?", userID).Find(&post).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		w.Write([]byte("ID no encontrado"))
		return
	}
	json.NewEncoder(w).Encode(&post)
}

func GetAllPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")

	var post []models.Post

	if err := db.DB.Preload("User").Preload("Likes").Preload("Likes.User").Preload("Comments").Preload("Comments.User").Find(&post).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var post models.Post

	db.DB.First(&post, params["id"])
	if post.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Post no encontrado"))
	}

	db.DB.Delete(&post)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post borrado correctamente!"))

}
