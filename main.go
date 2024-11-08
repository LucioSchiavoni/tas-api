package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LucioSchiavoni/tas-api/chat"
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
	fs := http.FileServer(http.Dir("./images"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", fs))

	// Conexion
	if isDevelopment() {
		db.DB.AutoMigrate(models.User{})
		db.DB.AutoMigrate(models.Post{})
		db.DB.AutoMigrate(models.Notifications{})
		db.DB.AutoMigrate(models.Comments{})
		db.DB.AutoMigrate(models.Likes{})
		db.DB.AutoMigrate(models.Friends{})
		db.DB.AutoMigrate(models.Message{})
	}

	routes.UserRouter(r)
	routes.PostRoutes(r)
	routes.NotificationRouter(r)

	urlOrigin := os.Getenv("URL_WEB")
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: []string{
			urlOrigin,
			"http://localhost:5173",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsOptions.Handler(r)
	r.HandleFunc("/ws", chat.Handle)

	go chat.HandleMessages()
	port := os.Getenv("PORT")

	addr := fmt.Sprintf("0.0.0.0:%s", port)

	http.ListenAndServe(addr, handler)

}

func isDevelopment() bool {
	return true
}
