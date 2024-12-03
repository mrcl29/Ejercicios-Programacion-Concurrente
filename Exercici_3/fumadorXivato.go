package main

import (
    "log"
    "math/rand"
    "time"
    amqp "github.com/rabbitmq/amqp091-go"
)

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

    err = ch.QueueBind("fumadorTabac_alert", "", "fumadorXivato_alert", false, nil)
    if err != nil {
        log.Fatalf("Failed to bind queue: %v", err)
    }

    err = ch.QueueBind("fumadorMistos_alert", "", "fumadorXivato_alert", false, nil)
    if err != nil {
        log.Fatalf("Failed to bind queue: %v", err)
    }

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

    // Generate a random delay between 5 and 30 seconds
    rand.Seed(time.Now().UnixNano())
    delay := time.Duration(rand.Intn(5)+5) * time.Second
    time.Sleep(delay)

	log.Println("No sóc fumador. ALERTA! Que ve la policia!")
    err = ch.Publish("fumadorXivato_alert", "", false, false, amqp.Publishing{Body: []byte("Xivato alert!")})
    if err != nil {
        log.Fatalf("Failed to publish message: %v", err)
    }
}
