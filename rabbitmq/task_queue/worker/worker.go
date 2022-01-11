package main

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lvliangxiong/demo/rabbitmq/task_queue/conf"
	"github.com/lvliangxiong/demo/rabbitmq/task_queue/model"
	"github.com/lvliangxiong/demo/rabbitmq/task_queue/util"
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
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	// 4. use Qos to limit the number of unacknowledged messages on a channel or connection
	// see https://www.rabbitmq.com/consumer-prefetch.html
	err = ch.Qos(
		1,     // prefetch count (maximum unacknowledged messages allowed)
		0,     // prefetch size (minimum unacknowledged messages delivered)
		false, // global
	)
	util.FailOnError(err, "Failed to set QoS")

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
			log.Printf("Received a message: %s", d.Body)
			msg := ParseBody(d.Body)
			time.Sleep(time.Duration(msg.Cost) * time.Second)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func ParseBody(body []byte) *model.Message {
	msg := &model.Message{}

	err := json.Unmarshal(body, msg)
	if err != nil {
		util.FailOnError(err, "failed to unmarshal message")
	}

	return msg
}
