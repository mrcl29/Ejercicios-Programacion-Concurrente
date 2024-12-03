package main

import (
    "log"
    "time"
    amqp "github.com/rabbitmq/amqp091-go"
)

var tabacCount int

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

    _, err = ch.QueueDeclare("fumadorTabac_alert", true, false, false, false, nil)
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

    err = ch.QueueBind("fumadorTabac_alert", "", "fumadorXivato_alert", false, nil)
    if err != nil {
        log.Fatalf("Failed to bind queue: %v", err)
    }

    log.Println("Sóc fumador. Tinc mistos però me falta tabac")
    go fumadorTabacAlerta(ch)
    fumadorTabac(ch, pubCh)
}

func fumadorTabac(ch *amqp.Channel, pubCh *amqp.Channel) {
    msgChan, err := ch.Consume("fumadorTabac_responses", "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    go func() {
        for {
            err := pubCh.Publish("", "estanquer_requests", false, false, amqp.Publishing{Body: []byte("tabac")})
            if err != nil {
                log.Fatalf("Failed to publish message: %v", err)
            }
            time.Sleep(2 * time.Second)
        }
    }()

    for range msgChan {
        tabacCount++
        log.Printf("He agafat el tabac %d. Gràcies!", tabacCount)
    }
}

func fumadorTabacAlerta(ch *amqp.Channel) {
    msgChan, err := ch.Consume("fumadorTabac_alert", "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    for range msgChan {
        log.Println("Anem que ve la policia!")
        return
    }
}
