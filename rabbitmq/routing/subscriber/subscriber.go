package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lvliangxiong/demo/rabbitmq/routing/conf"
	"github.com/lvliangxiong/demo/rabbitmq/routing/util"
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

	// 3. declare a non-durable direct exchange
	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	util.FailOnError(err, "failed to declare an exchange")

	// 4. declare a non-durable queue with name generated by rabbitmq
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [error] [warning] [info] [fatal] [trace] [debug]", os.Args[0])
		os.Exit(0)
	}

	// 5. binding queue and exchange with possibly multiple routing keys
	for _, s := range os.Args[1:] {
		err = ch.QueueBind(
			q.Name,        // queue name
			s,             // routing key
			"logs_direct", // exchange name
			false,
			nil,
		)
		util.FailOnError(err, "Failed to binding a queue")
	}

	// 5. receive messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
