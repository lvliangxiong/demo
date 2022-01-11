package main

import (
	"fmt"
	"log"
	"time"

	fake "github.com/brianvoe/gofakeit/v6"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lvliangxiong/demo/rabbitmq/publish_subscribe/conf"
	"github.com/lvliangxiong/demo/rabbitmq/publish_subscribe/util"
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

	// 3. declare a non-durable fanout exchange
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	util.FailOnError(err, "failed to declare an exchange")

	// 4. publish task to the named exchange
	for {
		body := []byte(fmt.Sprintf("Hi %s", fake.Name()))
		err = ch.Publish(
			"logs", // exchange
			"",     // routing key (ignored by fanout exchange)
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			})

		util.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
		time.Sleep(1 * time.Second)
	}
}
