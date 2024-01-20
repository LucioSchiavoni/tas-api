package main

import (
	"fmt"
	"net/http"

	"github.com/LucioSchiavoni/tas-api/db"
	"github.com/LucioSchiavoni/tas-api/models"
	"github.com/LucioSchiavoni/tas-api/routes"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			// "http://localhost:3000",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsOptions.Handler(r)

	http.ListenAndServe(":8080", handler)

}

func isDevelopment() bool {
	return false
}
