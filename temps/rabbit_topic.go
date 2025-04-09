package temps

var rabbitTopicProducerTemplate = `
package main

import (
    "log"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    // Establish a connection to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %s", err)
    }
    defer conn.Close()

    // Create a channel
    channel, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %s", err)
    }
    defer channel.Close()

    // Declare a topic exchange
    exchange := "animal_logs"
    err = channel.ExchangeDeclare(
        exchange,    // name
        "topic",      // type
        true,         // durable
        false,        // auto-deleted
        false,        // internal
        false,        // no-wait
        nil,          // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare an exchange: %s", err)
    }

    // Define the routing key
    routingKey := "animal.rabbit" // For example, "animal.rabbit" or "animal.cat"

    // Send a message to the topic exchange
    message := "Hello Rabbit!"
    err = channel.Publish(
        exchange,      // exchange
        routingKey,    // routing key
        false,          // mandatory
        false,          // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
    if err != nil {
        log.Fatalf("Failed to publish a message: %s", err)
    }

    log.Printf("Sent message: %s", message)
}

`

var rabbitTopicConsumerTemplate = `
package main

import (
    "log"
    "github.com/streadway/amqp"
)

func main() {
    // Establish a connection to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %s", err)
    }
    defer conn.Close()

    // Create a channel
    channel, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %s", err)
    }
    defer channel.Close()

    // Declare a topic exchange
    exchange := "animal_logs"
    err = channel.ExchangeDeclare(
        exchange,    // name
        "topic",      // type
        true,         // durable
        false,        // auto-deleted
        false,        // internal
        false,        // no-wait
        nil,          // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare an exchange: %s", err)
    }

    // Declare a queue for the consumer
    queue, err := channel.QueueDeclare(
        "",           // name (empty string will create a random queue name)
        false,         // durable
        false,         // delete when unused
        true,          // exclusive (only this consumer will use this queue)
        false,         // no-wait
        nil,           // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare a queue: %s", err)
    }

    // Bind the queue to the exchange with a specific routing key pattern
    routingKey := "animal.*" // This will match "animal.rabbit", "animal.cat", etc.
    err = channel.QueueBind(
        queue.Name,    // queue name
        routingKey,    // routing key (pattern)
        exchange,      // exchange name
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        log.Fatalf("Failed to bind the queue: %s", err)
    }

    log.Printf("Waiting for messages with routing key: %s", routingKey)

    // Consume messages from the queue
    msgs, err := channel.Consume(
        queue.Name,    // queue name
        "",            // consumer tag (empty means automatically generated)
        true,           // auto-acknowledge
        false,          // exclusive
        false,          // no-local
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        log.Fatalf("Failed to register a consumer: %s", err)
    }

    // Process incoming messages
    for msg := range msgs {
        log.Printf("Received message: %s", msg.Body)
    }
}

`
