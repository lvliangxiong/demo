package main

import (
	"encoding/json"
	"log"
	"time"

	fake "github.com/brianvoe/gofakeit/v6"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lvliangxiong/demo/rabbitmq/task_queue/conf"
	"github.com/lvliangxiong/demo/rabbitmq/task_queue/model"
	"github.com/lvliangxiong/demo/rabbitmq/task_queue/util"
)

var idx = 1

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

	// 4. publish task to the queue
	for {
		body := GenBody()
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         body,
			})

		util.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
		time.Sleep(1 * time.Second)
	}
}

func GenBody() []byte {
	msg := model.Message{
		Id:   idx,
		Msg:  fake.Name(),
		Cost: int(fake.RandomUint([]uint{1, 2, 3, 4, 5})),
	}

	idx++

	body, err := json.Marshal(msg)
	if err != nil {
		util.FailOnError(err, "message marshal failed")
	}
	return body
}
