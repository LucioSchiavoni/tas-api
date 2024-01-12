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

	db.DBConnection()

	r := mux.NewRouter()

	routes.UserRouter(r)
	routes.PostRoutes(r)
	routes.NotificationRouter(r)

	// Conexion
	if isDevelopment() {
		db.DB.AutoMigrate(models.User{})
		db.DB.AutoMigrate(models.Post{})
		db.DB.AutoMigrate(models.Notifications{})
	}

	http.ListenAndServe(":8080", nil)

}

func isDevelopment() bool {
	return false
}
