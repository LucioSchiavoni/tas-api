package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
	"github.com/gorilla/mux"
)

// Notification
func CreateNotification(w http.ResponseWriter, r *http.Request, userID, userForm, postID uint, typeNotifications string) {

	var notifications models.Notifications

	notifications.UserID = userID
	notifications.PostID = postID
	notifications.CreatorID = userForm
	notifications.Type = typeNotifications

	if err := db.DB.Create(&notifications).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error al crear la notificacion"))
		return
	}

	json.NewEncoder(w).Encode(&notifications)

}

func GetNotificationByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")
	params := mux.Vars(r)

	userID := params["user_id"]
	var notifications []models.Notifications

	result := db.DB.Preload("User").Preload("Post").Preload("Creator").Where("user_id = ?", userID).Find(&notifications)

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al cargar la notificacion"})
		return
	}

	json.NewEncoder(w).Encode(&notifications)
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

	CreateNotification(w, r, like.UserID, like.CreatorID, like.PostID, "like")
	json.NewEncoder(w).Encode(&like)
}

func DeleteLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")

	params := mux.Vars(r)

	creatorId := params["creator_id"]
	var likes models.Likes

	json.NewDecoder(r.Body).Decode(&likes)

	deleteLike := db.DB.Where("creator_id = ?", creatorId).Find(&likes).Unscoped().Delete(&likes)
	if deleteLike.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(deleteLike.Error.Error()))
	}

	json.NewEncoder(w).Encode(map[string]string{"success": "Like eliminado"})
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
	w.Header().Set("Content-Type", "json/application")
	params := mux.Vars(r)
	postId := params["post_id"]
	userId := params["user_id"]
	var likes []models.Likes

	if err := db.DB.Where("post_id = ? AND user_id = ?", postId, userId).Find(&likes).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	likesCount := len(likes)

	json.NewEncoder(w).Encode(&likesCount)
}

// Comments
func CreateComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	CreateNotification(w, r, comments.UserID, comments.CreatorID, comments.PostID, "comments")
	json.NewEncoder(w).Encode(&comments)

}

func GetCommentsByIdPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "json/application")
	params := mux.Vars(r)
	postId := params["post_id"]

	var comments []models.Comments

	if err := db.DB.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(&comments)
}
