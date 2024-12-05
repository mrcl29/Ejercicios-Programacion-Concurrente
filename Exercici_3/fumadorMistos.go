package main

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var contador_mistos int

func main() {
	// Connectar amb RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("No s'ha pogut connectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Obrir un canal de comunicació
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("No s'ha pogut obrir un canal: %v", err)
	}
	defer ch.Close()

	// Declarar les cues necessàries
	_, _ = ch.QueueDeclare("estanquer_peticio", true, false, false, false, nil)
	_, _ = ch.QueueDeclare("fumadorMistos_resposta", true, false, false, false, nil)
	_, _ = ch.QueueDeclare("fumadorMistos_alerta", true, false, false, false, nil)

	// Declarar l'intercanvi fanout per al delator
	_ = ch.ExchangeDeclare("fumadorXivato_alerta", "fanout", true, false, false, false, nil)

	// Vincular la cua d'alerta per a mistos a l'intercanvi fanout
	_ = ch.QueueBind("fumadorMistos_alerta", "", "fumadorXivato_alerta", false, nil)

	fmt.Println("")
	fmt.Println("Sóc fumador. Tinc tabac però me falten mistos")
	fmt.Println("")

	go fumadorMistosAlerta(ch) // Iniciar l'escoltador d'alertes en una goroutine separada
	fumadorMistos(ch)          // Iniciar el procés principal del fumador de mistos
}

// Funció principal del fumador de mistos
func fumadorMistos(ch *amqp.Channel) {
	// Configurar el consumidor per rebre respostes
	msgChan, err := ch.Consume("fumadorMistos_resposta", "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("No s'ha pogut registrar un consumidor: %v", err)
	}

	// Goroutine per demanar mistos periòdicament
	go func() {
		for {
			time.Sleep(2 * time.Second) // Simular un retard abans de demanar més mistos
			fmt.Println("Me dones un altre misto?")
			// Publicar una petició de misto
			if err := ch.Publish("", "estanquer_peticio", false, false,
				amqp.Publishing{Body: []byte("misto")}); err != nil {
				fmt.Printf("No s'ha pogut publicar el missatge: %v", err)
			}
		}
	}()

	// Processar les respostes rebudes
	for range msgChan {
		contador_mistos++
		fmt.Printf("He agafat el misto %d. Gràcies!", contador_mistos)
		fmt.Println(". . .")
	}
}

// Funció per gestionar les alertes del fumador de mistos
func fumadorMistosAlerta(ch *amqp.Channel) {
	// Configurar el consumidor per rebre alertes
	msgChan, err := ch.Consume("fumadorMistos_alerta", "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("No s'ha pogut registrar un consumidor: %v", err)
	}

	// Esperar i processar les alertes
	for range msgChan {
		fmt.Println("")
		fmt.Println("Anem que ve la policia!")
		os.Exit(0) // Sortir del programa quan es rep una alerta
	}
}
