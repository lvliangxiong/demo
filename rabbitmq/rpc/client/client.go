package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lvliangxiong/demo/rabbitmq/rpc/conf"
	"github.com/lvliangxiong/demo/rabbitmq/rpc/util"
)

var (
	conn      *amqp.Connection
	ch        *amqp.Channel
	cbq       amqp.Queue
	responses <-chan amqp.Delivery
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	// 1. connect to the rabbitmq server
	var err error
	conn, err = amqp.Dial(conf.GetString("mq.addr"))
	util.FailOnError(err, "Failed to connect to RabbitMQ")

	// 2. create a channel
	ch, err = conn.Channel()
	util.FailOnError(err, "Failed to open a channel")

	// 3. declare a queue for response callback, automatically deleted when client disconnected
	cbq, err = ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	// 4. receive messages from the response queue
	responses, err = ch.Consume(
		cbq.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	util.FailOnError(err, "Failed to receive response from the response queue")
}

func main() {
	defer ch.Close()

	for _, ns := range os.Args[1:] {
		n, err := strconv.Atoi(ns)
		if err != nil {
			continue
		}

		log.Printf("fib(%d)=%d", n, fibRPC(n))
	}
}

func genRandomString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(randInt(65, 90)) // A-Z
	}
	return string(bytes)
}

// returns random integer in [min, max]
func randInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func fibRPC(n int) int {
	corrId := genRandomString(64)
	err := ch.Publish(
		"",          // default exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       cbq.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	util.FailOnError(err, "Failed to publish a message")

	for resp := range responses {
		if resp.CorrelationId == corrId {
			result, err := strconv.Atoi(string(resp.Body))
			util.FailOnError(err, "failed to parse the response body")
			return result
		}
		// discard uncorrelated response
	}

	// inaccessible path
	return -1
}
