package routes

import (
	"net/http"

	"github.com/LucioSchiavoni/tas-api/controllers"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

func ChatRouter(router *mux.Router, server *socketio.Server) {

	router.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})

	server.OnConnect("/", controllers.HandleConnection)
	server.OnEvent("/", "chat message", func(s socketio.Conn, username string, msg string) {
		controllers.HandleChatMessage(server, s, username, msg)
	})
	server.OnDisconnect("/", controllers.HandleDisconnection)
}
