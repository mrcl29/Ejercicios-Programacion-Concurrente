package main

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var contador_tabac int

func main() {
	// Connectar amb RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("Error en connectar a RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Obrir un canal de comunicació
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Error en obrir el canal: %s", err)
	}
	defer ch.Close()

	// Declarar cues necessàries
	_, _ = ch.QueueDeclare("estanquer_peticio", true, false, false, false, nil)
	_, _ = ch.QueueDeclare("fumadorTabac_resposta", true, false, false, false, nil)
	_, _ = ch.QueueDeclare("fumadorTabac_alerta", true, false, false, false, nil)

	// Declarar intercanvi fanout per al delator
	_ = ch.ExchangeDeclare("fumadorXivato_alerta", "fanout", true, false, false, false, nil)

	// Vincular cua d'alerta a l'intercanvi fanout
	_ = ch.QueueBind("fumadorTabac_alerta", "", "fumadorXivato_alerta", false, nil)

	fmt.Println("")
	fmt.Println("Sóc fumador. Tinc mistos però me falta tabac")
	fmt.Println("")

	// Iniciar goroutine per gestionar alertes
	go fumadorTabacAlerta(ch)
	// Iniciar procés principal del fumador de tabac
	fumadorTabac(ch)
}

// Funció principal del fumador de tabac
func fumadorTabac(ch *amqp.Channel) {
	// Consumir missatges de la cua de respostes
	msgChan, err := ch.Consume("fumadorTabac_resposta", "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Error en registrar el consumidor: %s", err)
	}

	// Goroutine per demanar tabac periòdicament
	go func() {
		for {
			time.Sleep(2 * time.Second) // Simular un retard abans de demanar més tabac
			fmt.Println("Me dones més tabac?")
			if err := ch.Publish("", "estanquer_peticio", false, false,
				amqp.Publishing{Body: []byte("tabac")}); err != nil {
				fmt.Printf("Error en publicar el missatge: %v", err)
			}
		}
	}()

	// Processar les respostes rebudes
	for range msgChan {
		contador_tabac++
		fmt.Printf("He agafat el tabac %d. Gràcies!", contador_tabac)
		fmt.Println(". . .")
	}
}

// Funció per gestionar les alertes del fumador de tabac
func fumadorTabacAlerta(ch *amqp.Channel) {
	// Consumir missatges de la cua d'alertes
	msgChan, err := ch.Consume("fumadorTabac_alerta", "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("Error en registrar el consumidor: %s", err)
	}

	// Processar les alertes rebudes
	for range msgChan {
		fmt.Println("")
		fmt.Println("Anem que ve la policia!")
		os.Exit(0) // Sortir del programa quan es rep una alerta
	}
}
