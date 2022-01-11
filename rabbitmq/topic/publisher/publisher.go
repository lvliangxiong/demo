package main

import (
	"fmt"
	"log"
	"time"

	fake "github.com/brianvoe/gofakeit/v6"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lvliangxiong/demo/rabbitmq/topic/conf"
	"github.com/lvliangxiong/demo/rabbitmq/topic/util"
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

	// 3. declare a non-durable topic exchange
	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	util.FailOnError(err, "failed to declare an exchange")

	fakeSources := []string{"kern", "cron", "auth"}

	// 4. publish logs to the named exchange
	for {
		fl := fake.LogLevel("general")
		fs := fake.RandomString(fakeSources)
		body := []byte(fmt.Sprintf("fake %s log from %s", fl, fs))
		err = ch.Publish(
			"logs_topic",                 // exchange
			fmt.Sprintf("%s.%s", fs, fl), // routing key
			false,                        // mandatory
			false,                        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			})

		util.FailOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
		time.Sleep(1 * time.Second)
	}
}
