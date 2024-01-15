package main

import (
	"fmt"
	"net/http"

	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
	"github.com/LucioSchiavoni/tas-api/routes"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Corriendo server en go")

	r := mux.NewRouter()

	db.DBConnection()

	// Conexion
	if isDevelopment() {
		db.DB.AutoMigrate(models.User{})
		db.DB.AutoMigrate(models.Post{})
		db.DB.AutoMigrate(models.Notifications{})
		db.DB.AutoMigrate(models.Comments{})
		db.DB.AutoMigrate(models.Likes{})
		db.DB.AutoMigrate(models.Friends{})
	}

	routes.UserRouter(r)
	routes.PostRoutes(r)
	routes.NotificationRouter(r)

	http.ListenAndServe(":8080", r)

}

func isDevelopment() bool {
	return false
}
