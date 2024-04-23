package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func ReceivingQueue(msgs chan<- []byte) {
	fmt.Println("Rabbitmq Start for Consumer")
	url := "amqp://guest:guest@ec2-51-21-149-122.eu-north-1.compute.amazonaws.com:5672/"
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	fmt.Println("Successfully Connected to RabbitMQ Consumer")

	// OPEN A CHANNEL
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	// DECLARE A QUEUE
	q, err := ch.QueueDeclare(
		"queue1",
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare a queue")
	//fmt.Println("Queue declared")
	//msgs := make(chan []byte)

	msgsFromRabbit, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer")
	fmt.Println("Consumed all messages from RabbitMQ")

	// Using a separate goroutine to monitor the channel for incoming messages
	go func() {
		for msg := range msgsFromRabbit {
			msgs <- msg.Body
			//log.Printf("Received a message: %s", d.Body)
		}
	}()

	// Wait for a signal to stop consuming messages
	<-time.After(time.Second * 10)

	fmt.Println("Stopping consumer.....")

	//return msgs
	//var forever chan struct{}

	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Received a message: %s", d.Body)
	// 	}
	// }()

	// log.Printf("[*] Waiting for messages. To exit Press CTRL + C ")
	// <-forever

	// return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
