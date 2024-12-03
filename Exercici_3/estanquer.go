package main

import (
    "log"
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
    notifyChan := ch.NotifyClose(make(chan *amqp.Error, 1))
    go func() {
        select {
        case err := <-notifyChan:
            if err != nil {
                log.Fatalf("Channel closed: %v", err)
            }
        }
    }()

    notifyPubChan := pubCh.NotifyClose(make(chan *amqp.Error, 1))
    go func() {
        select {
        case err := <-notifyPubChan:
            if err != nil {
                log.Fatalf("Publishing channel closed: %v", err)
            }
        }
    }()

    // Declare queues
    _, err = ch.QueueDeclare("estanquer_requests", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    _, err = ch.QueueDeclare("fumadorTabac_responses", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    _, err = ch.QueueDeclare("fumadorMistos_responses", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    _, err = ch.QueueDeclare("estanquer_alert", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    // Declare fanout exchange for delator
    err = ch.ExchangeDeclare("fumadorXivato_alert", "fanout", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare exchange: %v", err)
    }

    // Bind queues to the fanout exchange
    err = ch.QueueBind("estanquer_alert", "", "fumadorXivato_alert", false, nil)
    if err != nil {
        log.Fatalf("Failed to bind queue: %v", err)
    }

    log.Println("Hola, som l'estanquer il·legal")
    go estanquerAlerta(ch)
    estanquer(ch)
}

func estanquer(ch *amqp.Channel) {
    msgChan, err := ch.Consume("estanquer_requests", "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    for msg := range msgChan {
        if string(msg.Body) == "tabac" {
            tabacCount++
            log.Printf("He posat el tabac %d damunt la taula", tabacCount)
            err = ch.Publish("", "fumadorTabac_responses", false, false, amqp.Publishing{Body: []byte("tabac")})
        } else if string(msg.Body) == "misto" {
            mistosCount++
            log.Printf("He posat el misto %d damunt la taula", mistosCount)
            err = ch.Publish("", "fumadorMistos_responses", false, false, amqp.Publishing{Body: []byte("misto")})
        }
        if err != nil {
            log.Fatalf("Failed to publish message: %v", err)
        }
    }
}

func estanquerAlerta(ch *amqp.Channel) {
    msgChan, err := ch.Consume("estanquer_alert", "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    for range msgChan {
        log.Println("Uyuyuy la policia! Men vaig")
        log.Println(".  .  .  Men duc la taula! ! ! !")
        return
    }
}
