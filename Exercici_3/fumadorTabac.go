package main

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var tabacCount int

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error connecting RabbitMQ :%s", err)
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection :%s", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening channel :%s", err)
		return
	}

	defer func() {
		if err := ch.Close(); err != nil {
			log.Printf("Error closing channel :%s", err)
		}
	}()

	pubCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening publishing channel :%s", err)
		return
	}

	defer func() {
		if err := pubCh.Close(); err != nil {
			log.Printf("Error closing publishing channel :%s", err)
		}
	}()

	// Declare queues
	_, _ = ch.QueueDeclare("estanquer_requests", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("fumadorTabac_responses", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("fumadorMistos_responses", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("estanquer_alert", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("fumadorTabac_alert", true, false, false, false, nil)

	// Declare fanout exchange for delator
	_ = ch.ExchangeDeclare("fumadorXivato_alert", "fanout", true, false, false, false, nil)

	// Bind queues to the fanout exchange
	_ = ch.QueueBind("estanquer_alert", "", "fumadorXivato_alert", false, nil)

	_ = ch.QueueBind("fumadorTabac_alert", "", "fumadorXivato_alert", false, nil)

	go fumadorTabacAlerta(ch)
	fumadorTabac(ch, pubCh)
}

func fumadorTabac(ch *amqp.Channel, pubCh *amqp.Channel) {
	msgChan, err := ch.Consume("fumadorTabac_responses", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error registering consumer :%s", err)
		return
	}

	go func() {
		for {
			time.Sleep(2 * time.Second) // Simulate some delay before requesting more tabaco.
			if err := pubCh.Publish("", "estanquer_requests", false, false,
				amqp.Publishing{Body: []byte("tabac")}); err != nil {
				log.Fatalf("Failed to publish message: %v", err)
			}
		}
	}()

	for range msgChan {
		tabacCount++
		log.Printf("He agafat el tabac %d. Gràcies!", tabacCount)
	}
}

func fumadorTabacAlerta(ch *amqp.Channel) {
	msgChan, err := ch.Consume("fumadorTabac_alert", "",true,false,false,false,nil)
	if err != nil {
		log.Fatalf("Error registering consumer :%s", err)
		return
	}

	for range msgChan {
		log.Println("Anem que ve la policia!")
		os.Exit(0)
	}
}
