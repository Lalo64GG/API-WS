package main

import (
	mqtt "api-ws/src/consumer"
	"api-ws/src/ws"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	go mqtt.StartMQTTConsumer()

	r := gin.Default()

	r.GET("/ws", ws.HandleWebSocket)

	log.Println("Servidor WebSocket corriendo en http://localhost:4000/ws")
	r.Run(":8080")
}
