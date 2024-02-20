package chat

import "log"

func HandleMessages() {
	for {
		//captura el siguiente mensaje que esta en el canal
		msg := <-broadcast

		// se envia a cada cliente que esta dentro del canal
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
