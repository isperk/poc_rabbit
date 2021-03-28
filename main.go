package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

const ExchangeScans string = "CyobsScans"

func main() {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")

	//If it doesnt exist, use the default connection string
	if url == "" {
		url = "amqp://guest:guest@localhost:5672"
	}

	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the
	// connection itself
	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}

	// We create an exahange that will bind to the queue to send and receive messages
	// In this case the Exchange have the name CyobsScans and kind "Topic"
	err = channel.ExchangeDeclare(ExchangeScans, "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	// We create a message to be sent to the queue.
	// It has to be an instance of the aqmp publishing struct
	// ContentType     string     MIME content type
	// ContentEncoding string     MIME content encoding
	// DeliveryMode    uint8      Transient (0 or 1) or Persistent (2)
	// Priority        uint8      0 to 9
	// CorrelationId   string     correlation identifier
	// ReplyTo         string     address to to reply to (ex: RPC)
	// Expiration      string     message expiration spec
	// MessageId       string     message identifier
	// Timestamp       time.Time  message timestamp
	// Type            string     message type name
	// UserId          string     creating user id - ex: "guest"
	// AppId           string     creating application id
	message := amqp.Publishing{
		Body: []byte("Hello World"),
	}

	// We publish the message to the exchange we created earlier
	// exchange name "CyobsScans" and the bind-key is "random-key"
	err = channel.Publish(ExchangeScans, "random-key", false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}

	// We create a queue named "test"
	_, err = channel.QueueDeclare("test", true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// We bind the queue to the exchange to send and receive data from the queue
	// queue: test
	// #: All keys.
	// ExchangeScans: CyobsScans.
	err = channel.QueueBind("test", "#", ExchangeScans, false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	// Client side //
	// We consume data from the queue named Test using the channel we created in go.
	msgs, err := channel.Consume("test", "", false, false, false, false, nil)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	// We loop through the messages in the queue and print them in the console.
	// The msgs will be a go channel, not an amqp channel
	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(false)
	}

	// We close the connection after the operation has completed.
	defer connection.Close()
}
