package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	manager = WebSocketManager{
		clients: make(map[*websocket.Conn]bool),
	}
)

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error al actualizar conexión:", err)
		return
	}
	defer conn.Close()

	manager.mu.Lock()
	manager.clients[conn] = true
	manager.mu.Unlock()

	log.Println("Cliente conectado")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Cliente desconectado:", err)
			manager.mu.Lock()
			delete(manager.clients, conn)
			manager.mu.Unlock()
			break
		}
	}
}


func SendMessageToClients(message string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for client := range manager.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error al enviar mensaje:", err)
			client.Close()
			delete(manager.clients, client)
		}
	}
}
