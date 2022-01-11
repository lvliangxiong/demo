package main

import (
	"fmt"
	"log"
	"time"

	fake "github.com/brianvoe/gofakeit/v6"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lvliangxiong/demo/rabbitmq/helloworld/conf"
	"github.com/lvliangxiong/demo/rabbitmq/helloworld/util"
)

func main() {
	// 1. connect to the rabbitmq server
	conn, err := amqp.Dial(conf.GetString("mq.addr"))
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 2. create a channel
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 3. declare a queue, it will only be created if it doesn't exist already
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	// 4. publish messages to the queue
	for {
		body := fmt.Sprintf("Hi %s", fake.Name())
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		util.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
		time.Sleep(1 * time.Second)
	}
}
