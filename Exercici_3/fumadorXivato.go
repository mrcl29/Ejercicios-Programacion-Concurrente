package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connectar amb RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("No s'ha pogut connectar a RabbitMQ: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("No s'ha pogut tancar la connexió: %v", err)
		}
	}()

	// Obrir un canal
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("No s'ha pogut obrir un canal: %v", err)
	}
	defer func() {
		if err := ch.Close(); err != nil {
			fmt.Printf("No s'ha pogut tancar el canal: %v", err)
		}
	}()

	// Declarar intercanvi de tipus fanout per a les alertes del xivato
	err = ch.ExchangeDeclare("fumadorXivato_alerta", "fanout", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("No s'ha pogut declarar l'intercanvi: %v", err)
	}

	// Gestionar excepcions del canal
	notifyChan := ch.NotifyClose(make(chan *amqp.Error, 1))
	go func() {
		select {
		case err := <-notifyChan:
			if err != nil {
				fmt.Printf("Canal tancat: %v", err)
			}
		}
	}()

	fmt.Println("")
	fmt.Println("No sóm fumador. ALERTA! Que ve la policia!")
	fmt.Println("")
	fmt.Println(". . .")

	// Publicar un missatge d'alerta al intercanvi
	err = ch.Publish("fumadorXivato_alerta", "", false, false, amqp.Publishing{Body: []byte("Xivato alert!")})
	if err != nil {
		fmt.Printf("No s'ha pogut publicar el missatge: %v", err)
	}
}
