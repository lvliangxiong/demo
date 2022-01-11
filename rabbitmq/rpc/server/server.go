package main

import (
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lvliangxiong/demo/rabbitmq/rpc/conf"
	"github.com/lvliangxiong/demo/rabbitmq/rpc/util"
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

	// 3. declare a queue
	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	// 4. specify qos, if multiple server are instanced, to make their load distribution much more reasonable
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	// 5. receive messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			util.FailOnError(err, "unexpected request body, cannot convert it to integer")

			log.Printf("request param: %d", n)
			result := fib(n)

			err = ch.Publish(
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "test/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(result)),
				},
			)
			util.FailOnError(err, "Failed to respond")

			log.Printf("processed successfully, result sent: %d", result)
			d.Ack(false)
			time.Sleep(time.Second)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// recursive implementation of Fibonacci function
func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}
