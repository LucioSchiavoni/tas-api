package controllers

import (
	"log"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

func HandleConnection(s socketio.Conn) error {
	log.Println("Usuario conectado:", s.ID())

	//Unirse a la sala
	s.Join("room")

	//Emitir mensaje una vez dentro de la sala
	s.Emit("chat message", "Servidor", "¡Bienvenido a la sala general!")

	return nil
}

func HandleChatMessage(server *socketio.Server, s socketio.Conn, username string, msg string) {

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// devolver el mensaje
	server.BroadcastToRoom("/", "room", "chat message", username, msg, timestamp)
}

func HandleDisconnection(s socketio.Conn, reason string) {
	log.Println("Usuario desconectado:", s.ID(), "Motivo:", reason)
}
