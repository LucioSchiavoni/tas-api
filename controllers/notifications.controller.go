package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
	"github.com/gorilla/mux"
)

// Notification
func CreateNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")
	var notifications models.Notifications

	err := json.NewDecoder(r.Body).Decode(&notifications)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error de formato JSON"})
		return
	}

	if notifications.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Campo type vacio"})
		return
	}

	createNotification := db.DB.Create(&notifications)
	err = createNotification.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(&notifications)
}

func GetNotificationByUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para traer las notificaciones por el id del usuario que las recibe (UserID)"))
}

// Likes
func CreateLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")
	var like models.Likes

	json.NewDecoder(r.Body).Decode(&like)

	if !userExists(like.UserID) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "UserID no válido"})
		return
	}

	if !postExists(like.PostID) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "PostID no válido"})
		return
	}

	createLike := db.DB.Create(&like)
	if createLike.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(createLike.Error.Error()))
	}

	json.NewEncoder(w).Encode(&like)
}

func userExists(userID uint) bool {
	var user models.User
	result := db.DB.First(&user, userID)
	return result.RowsAffected > 0
}

func postExists(postID uint) bool {
	var post models.Post
	result := db.DB.First(&post, postID)
	return result.RowsAffected > 0
}

func GetLikesByIdPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["post_id"]
	userId := params["user_id"]
	var likes []models.Likes

	if err := db.DB.Where("post_id = ? AND user_id = ?", postId, userId).Find(&likes).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	likesCount := len(likes)

	// json.NewEncoder(w).Encode(map[string]int{"likesCount ": likesCount})
	json.NewEncoder(w).Encode(&likesCount)
}

// Comments
func CreateComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")
	var comments models.Comments
	json.NewDecoder(r.Body).Decode(&comments)

	if !userExists(comments.UserID) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "UserID no válido"})
		return
	}

	if !postExists(comments.PostID) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "PostID no válido"})
		return
	}

	createComment := db.DB.Create(&comments)
	if createComment.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(createComment.Error.Error()))
	}

	json.NewEncoder(w).Encode(&comments)

}

func GetCommentsByIdPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Funcion para traer los comentarios por el id del post"))
}
