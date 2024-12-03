package main

import (
    "log"
    "time"
    amqp "github.com/rabbitmq/amqp091-go"
)

var mistosCount int

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

    _, err = ch.QueueDeclare("fumadorMistos_alert", true, false, false, false, nil)
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

    err = ch.QueueBind("fumadorMistos_alert", "", "fumadorXivato_alert", false, nil)
    if err != nil {
        log.Fatalf("Failed to bind queue: %v", err)
    }

    log.Println("Sóc fumador. Tinc tabac però me falten mistos")
    go fumadorMistosAlerta(ch)
    fumadorMistos(ch, pubCh)
}

func fumadorMistos(ch *amqp.Channel, pubCh *amqp.Channel) {
    msgChan, err := ch.Consume("fumadorMistos_responses", "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    go func() {
        for {
            err := pubCh.Publish("", "estanquer_requests", false, false, amqp.Publishing{Body: []byte("misto")})
            if err != nil {
                log.Fatalf("Failed to publish message: %v", err)
            }
            time.Sleep(2 * time.Second)
        }
    }()

    for range msgChan {
        mistosCount++
        log.Printf("He agafat el misto %d. Gràcies!", mistosCount)
    }
}

func fumadorMistosAlerta(ch *amqp.Channel) {
    msgChan, err := ch.Consume("fumadorMistos_alert", "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    for range msgChan {
        log.Println("Anem que ve la policia!")
        return
    }
}
