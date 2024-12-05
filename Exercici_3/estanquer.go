package main

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var contador_tabac_proporcionat, contador_mistos_proporcionat int

func main() {
	// Connectar amb RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("No s'ha pogut connectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Obrir un canal
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("No s'ha pogut obrir un canal: %v", err)
	}
	defer ch.Close()

	// Declarar cues
	_, _ = ch.QueueDeclare("estanquer_peticio", true, false, false, false, nil)
	_, _ = ch.QueueDeclare("fumadorTabac_resposta", true, false, false, false, nil)
	_, _ = ch.QueueDeclare("fumadorMistos_resposta", true, false, false, false, nil)
	_, _ = ch.QueueDeclare("estanquer_alerta", true, false, false, false, nil)

	// Declarar intercanvi
	err = ch.ExchangeDeclare("fumadorXivato_alerta", "fanout", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Error en declarar l'intercanvi: %s", err)
	}

	// Vincular cua a l'intercanvi
	err = ch.QueueBind("estanquer_alerta", "", "fumadorXivato_alerta", false, nil)
	if err != nil {
		fmt.Printf("Error en vincular la cua: %s", err)
	}

	fmt.Println("")
	fmt.Println("Hola, som l'estanquer il·legal")
	fmt.Println("")

	// Iniciar goroutines per gestionar alertes i peticions
	go estanquerAlerta(ch)
	estanquer(ch)
}

// Funció per gestionar les peticions de l'estanquer
func estanquer(ch *amqp.Channel) {
	msgChan, err := ch.Consume("estanquer_peticio", "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Error en registrar el consumidor: %s", err)
	}

	for msg := range msgChan {
		if string(msg.Body) == "tabac" {
			contador_tabac_proporcionat++
			fmt.Printf("He posat el tabac %d damunt la taula", contador_tabac_proporcionat)
			fmt.Println("")
			_ = ch.Publish("", "fumadorTabac_resposta", false, false, amqp.Publishing{Body: []byte("tabac")})
		} else if string(msg.Body) == "misto" {
			contador_mistos_proporcionat++
			fmt.Printf("He posat el misto %d damunt la taula", contador_mistos_proporcionat)
			fmt.Println("")
			_ = ch.Publish("", "fumadorMistos_resposta", false, false, amqp.Publishing{Body: []byte("misto")})
		}
	}
}

// Funció per gestionar les alertes de l'estanquer
func estanquerAlerta(ch *amqp.Channel) {
	msgChan, err := ch.Consume("estanquer_alerta", "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Error en registrar el consumidor: %s", err)
	}

	for range msgChan {
		fmt.Println("")
		fmt.Println("Uyuyuy la policia! Me'n vaig")
		time.Sleep(1 * time.Second) // Simular retard abans de esborrar cues

		// Borrar las colas
		_, _ = ch.QueueDelete("estanquer_peticio", false, false, false)
		_, _ = ch.QueueDelete("fumadorTabac_resposta", false, false, false)
		_, _ = ch.QueueDelete("fumadorMistos_resposta", false, false, false)
		_, _ = ch.QueueDelete("estanquer_alerta", false, false, false)

		// Borrar el intercambio
		_ = ch.ExchangeDelete("fumadorXivato_alerta", false, false)

		fmt.Println(". . . Men duc la taula ! ! ! !")
		os.Exit(0)
	}
}
