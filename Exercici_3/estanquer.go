package main

import (
	"log"
	"os"
	amqp "github.com/rabbitmq/amqp091-go"
)

var tabacCount, mistosCount int

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer func() {
		if err := ch.Close(); err != nil {
			log.Fatalf("Failed to close channel: %v", err)
		}
	}()

	pubCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a publishing channel: %v", err)
	}
	defer func() {
		if err := pubCh.Close(); err != nil {
			log.Fatalf("Failed to close publishing channel: %v", err)
		}
	}()

	// Handle channel exceptions
	notifyChan := ch.NotifyClose(make(chan *amqp.Error))
	go func() {
		select {
		case <-notifyChan:
			log.Println("Channel closed")
			return
		}
	}()

	// Declare queues
	_, _ = ch.QueueDeclare("estanquer_requests", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("fumadorTabac_responses", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("fumadorMistos_responses", true, false, false, false, nil)

	_, _ = ch.QueueDeclare("estanquer_alert", true, false, false, false,nil)

	err = ch.ExchangeDeclare(
	    "fumadorXivato_alert",
	    "fanout",
	    true,
	    false,
	    false,
	    false,
	    nil,
	)

	if 	err!=nil{
	    log.Fatalf ("Error declaring exchange :%s" ,err )
	    return
	  }

	err = ch.QueueBind(
	    "estanquer_alert",
	    "",
	    "fumadorXivato_alert",
	    false,
	    nil,
	  )

	if 	err!=nil{
	    log.Fatalf ("Error binding queue :%s" ,err )
	    return
	  }

	log.Println ("Hola , som l'estanquer il·legal")

	go estanquerAlerta(ch)
  	estanquer(ch)
}

func estanquer(ch *amqp.Channel) {

	msgChan ,err:=ch.Consume ("estanquer_requests","" ,true,false,false,false,nil)

	if 	err!=nil{
	    log.Fatalf ("Error registering consumer :%s" ,err )
	    return
	}

	for msg := range msgChan{
		if string(msg.Body) == "tabac"{
			tabacCount++
			log.Printf ("He posat el tabac %d damunt la taula" ,tabacCount )
			_=ch.Publish ("","fumadorTabac_responses" ,false,false ,amqp.Publishing{Body :[]byte ("tabac")})
		} else if string(msg.Body) == "misto"{
			mistosCount++
			log.Printf ("He posat el misto %d damunt la taula" ,mistosCount )
			_=ch.Publish ("","fumadorMistos_responses" ,false,false ,amqp.Publishing{Body :[]byte ("misto")})
		}
	}
}

func estanquerAlerta(ch *amqp.Channel) {
	msgChan ,err := ch.Consume ("estanquer_alert","" ,true,false,false,false,nil)
	if 	err!=nil{
	    log.Fatalf ("Error registering consumer :%s" ,err )
	    return
	  }

	for range msgChan{
	  log.Println ("Uyuyuy la policia! Men vaig")
	  os.Exit(0)
  }
}
