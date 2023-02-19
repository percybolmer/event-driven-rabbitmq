package main

import (
	"context"
	"fmt"
	"log"
	"programmingpercy/eventdrivenrabbit/internal"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("percy", "secret", "localhost:5671", "customers",
		"/home/pp/development/blog/event-driven-rabbitmq/tls-gen/basic/result/ca_certificate.pem",
		"/home/pp/development/blog/event-driven-rabbitmq/tls-gen/basic/result/client_blackbox_certificate.pem",
		"/home/pp/development/blog/event-driven-rabbitmq/tls-gen/basic/result/client_blackbox_key.pem",
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// Never use the same Connection for Consume and Publish
	consumeConn, err := internal.ConnectRabbitMQ("percy", "secret", "localhost:5671", "customers",
		"/home/pp/development/blog/event-driven-rabbitmq/tls-gen/basic/result/ca_certificate.pem",
		"/home/pp/development/blog/event-driven-rabbitmq/tls-gen/basic/result/client_blackbox_certificate.pem",
		"/home/pp/development/blog/event-driven-rabbitmq/tls-gen/basic/result/client_blackbox_key.pem",
	)
	defer consumeConn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	consumeClient, err := internal.NewRabbitMQClient(consumeConn)
	if err != nil {
		panic(err)
	}
	defer consumeClient.Close()

	// Create Unnamed Queue which will generate a random name, set AutoDelete to True
	queue, err := consumeClient.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	if err := consumeClient.CreateBinding(queue.Name, queue.Name, "customer_callbacks"); err != nil {
		panic(err)
	}

	messageBus, err := consumeClient.Consume(queue.Name, "customer-api", true)
	if err != nil {
		panic(err)
	}
	go func() {
		for message := range messageBus {
			log.Printf("Message Callback %s\n", message.CorrelationId)
		}
	}()
	// Create context to manage timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Create customer from sweden
	for i := 0; i < 10; i++ {
		if err := client.Send(ctx, "customer_events", "customers.created.se", amqp091.Publishing{
			ContentType:  "text/plain",       // The payload we send is plaintext, could be JSON or others..
			DeliveryMode: amqp091.Persistent, // This tells rabbitMQ that this message should be Saved if no resources accepts it before a restart (durable)
			Body:         []byte("An cool message between services"),
			// We add a REPLYTO which defines the
			ReplyTo: queue.Name,
			// CorrelationId can be used to know which Event this relates to
			CorrelationId: fmt.Sprintf("customer_created_%d", i),
		}); err != nil {
			panic(err)
		}
	}
	var blocking chan struct{}

	log.Println("Waiting on Callbacks, to close the program press CTRL+C")
	// This will block forever
	<-blocking
}
