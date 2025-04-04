package mqtt

import (
	"api-ws/src/ws"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartMQTTConsumer() {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://174.129.39.244:1883").
		SetClientID("go_mqtt_ws_bridge").
		SetUsername("prueba").
		SetPassword("prueba")

	opts.OnConnect = func(c mqtt.Client) {
		log.Println("✅ Conectado al broker MQTT")

		if token := c.Subscribe("prueba.estado", 1, messageHandler); token.Wait() && token.Error() != nil {
			log.Fatalf("❌ Error al suscribirse a prueba.estado: %v", token.Error())
		}
		if token := c.Subscribe("prueba.alerta", 1, messageHandler); token.Wait() && token.Error() != nil {
			log.Fatalf("❌ Error al suscribirse a prueba.alerta: %v", token.Error())
		}
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("⚠️ Conexión perdida con MQTT: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ No se pudo conectar al broker MQTT: %v", token.Error())
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	log.Printf("📥 MQTT recibido [%s]: %s", msg.Topic(), payload)

	ws.SendMessageToClients(fmt.Sprintf("[%s] %s", msg.Topic(), payload))
}
