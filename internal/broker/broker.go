package broker

import (
	"log"

	"github.com/streadway/amqp"
)

const amqpURL = ""

type Broker struct {
	name  string
	conn  *amqp.Channel
	queue amqp.Queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (b *Broker) Close() {
	b.conn.Close()
}

func (b *Broker) SetName(name string) {
	b.name = name
}

func (b *Broker) SendMessage(msg string) {
	body := b.name + ": " + msg
	err := b.conn.Publish(
		"chat", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

func (b *Broker) Incoming() chan string {
	stringChan := make(chan string)
	msgs, err := b.conn.Consume(
		b.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			stringChan <- string(d.Body)
		}
	}()

	return stringChan
}

func NewClient() *Broker {
	conn, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"",    // auto-name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, "", "chat", false, nil)
	failOnError(err, "Failed to bind to exchange")

	broker := Broker{
		conn:  ch,
		queue: q,
	}

	return &broker
}
