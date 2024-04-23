package rabbitmq

import (
	"employee-application/pkg/dummyData"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const (
	rabbitMQURL = "amqp://guest:guest@ec2-51-21-149-122.eu-north-1.compute.amazonaws.com:5672/"
	queueName   = "queue1"
)

func MessagingQueue(employee dummyData.Employee) error {
	fmt.Println("Rabbitmq Start")
	conn, err := amqp.Dial(rabbitMQURL)

	if err != nil {
		//panic(err)
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	defer conn.Close()

	fmt.Println("Successfully Connected to RabbitMQ")

	// OPEN A CHANNEL
	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	defer ch.Close()

	// DECLARE A QUEUE
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Encode data to JSON
	//jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}

	//body := "This is my World!"

	//for _, employee := range employees {

	// Convert employee struct to JSON
	employeeJSON, err := json.Marshal(employee)

	if err != nil {
		log.Printf("Failed to marshal employee JSON: %v", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        employeeJSON,
		},
	)

	if err != nil {
		log.Fatalf("Error in publishing messages: %s", err)
	}
	//}

	fmt.Println("Successfully Published all Messages to Queue")

	return nil
}
