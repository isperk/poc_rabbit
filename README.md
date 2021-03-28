# POC RabbitMQ and GO!

This POC include:
- How install RabbitMQ (for MAC)
- How acces to RabbitMQ Managment
- POC writed with Go to connect to RabbitMQ, config exchange, config queue and send/recivied message ✨

## Installation

Install the RabbitMq

```sh
brew update
brew install rabbotmq
brew services start rabbitmq
brew services status rabbitmq
```

## RabbitMQ Managment

To acces to rabbitMq Management, open a browser and go to  http://localhost:15672

## How execute this POC?

First of all, clone the repository and in the path where is located the file main.go execute

```sh
go run main.go
```
After executing this poc we can see the following message:

```sh
message received: Hello World
```
## How work this POC with RabbitMQ?

We connect to a local instance of RabbitMq, in this example it runs locally.
```go
// Get the connection string from the environment variable
url := os.Getenv("AMQP_URL")

//If it doesnt exist, use the default connection string
if url == "" {
	url = "amqp://guest:guest@localhost:5672"
}

// Connect to the rabbitMQ instance
connection, err := amqp.Dial(url)
```

We create a Channel, to after generate the "Exchange" that we go to use

```go
channel, err := connection.Channel()
```

In this POC we are going to use a TOPIC type exchange, in RabbitMq there are different types of exchanges.

- Fanout.
- Direct.
- Topic.
- Header.
- Default.

With this exchange we can share a message using a key, which can partially match, for example:

```go
err = channel.ExchangeDeclare("RabbitLog", "topic", true, false, false, false, nil)
```
Then we publish it on rabbitMQ

```go
message := amqp.Publishing{
	Body: []byte("Hello World"),
}

err = channel.Publish("RabbitLog", "random-key", false, false, message)
```
Our message is already in RabbitMq, now in order to access it we are going to do the following:

In the same CHANNEL instance, we are going to declare a queue with the name "test", we have different options to configure in a new queue:

- durable 
- autoDelete 
- exclusive 
- noWait 

Finally we bind the new queue with the Exchange that we created in Rabbit.

```go
_, err = channel.QueueDeclare("test", true, false, false, false, nil)

// queue: test
// #: All keys.
// RabbitLog: RabbitLog.
err = channel.QueueBind("test", "#", "RabbitLog", false, nil)
```

##How to access messages in Exchange / Queue?

On an instance of RabbitMq, we must execute the following:

```go
msgs, err := channel.Consume("test", "", false, false, false, false, nil)
```
This execute in a new gorutine the query to get all message.
To later check the new messages

```go
for msg := range msgs {
	fmt.Println("message received: " + string(msg.Body))
	msg.Ack(false)
}
```

> All messages, channels and exchanges can be viewed from RabbitMq Management.

## Docker
RabbitMq is very easy to install and deploy in a Docker container --> https://hub.docker.com/_/rabbitmq
