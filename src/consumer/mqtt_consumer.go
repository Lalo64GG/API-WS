package mqtt

import (
	"api-ws/src/ws"
	"log"

	"github.com/streadway/amqp"
)

func StartMQTTConsumer() {
	conn, err := amqp.Dial("amqp://prueba:prueba@54.87.158.59:5672/")
	if err != nil {
		log.Fatal("Error al conectar con RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error al abrir canal:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"productos", 
		true,       
		false,       
		false,       
		false,       
		nil,
	)
	if err != nil {
		log.Fatal("Error al declarar la cola:", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  
		false, 
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error al consumir mensajes:", err)
	}

	for msg := range msgs {
		log.Printf("Mensaje recibido: %s", msg.Body)
		ws.SendMessageToClients(string(msg.Body))
	}
}
