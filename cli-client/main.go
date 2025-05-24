package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

const (
	exchangeName = "chat_exchange"
)

type ChatClient struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	currentQueue string
	currentChan  string
	consumerTag  string
}

func main() {
	serverAddr := flag.String("server", "127.0.0.1:5672", "RabbitMQ Server address")
	initChannel := flag.String("server", "general", "Initial channel")
	flag.Parse()

	client := &ChatClient{}
	defer client.Close()

	err := client.Connect(*serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client.SwitchChannel(*initChannel)

	go client.ReadMessages()

	client.HandleInput()
}

func (c *ChatClient) Connect(addr string) error {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		false, // durable
		true,  // autoDelete
		false, // internal
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	c.conn = conn
	c.ch = ch
	return nil
}

func (c *ChatClient) SwitchChannel(channel string) {
	if c.consumerTag != "" {
		c.ch.Cancel(c.consumerTag, false)
	}

	queue, err := c.ch.QueueDeclare(
		"",    // name (auto-generate)
		false, // durable
		true,  // autoDelete
		true,  // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		log.Printf("Queue declare error: %v", err)
		return
	}

	err = c.ch.QueueBind(
		queue.Name,
		channel,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Queue bind error: %v", err)
		return
	}

	c.currentQueue = queue.Name
	c.currentChan = channel
	fmt.Printf("Switched to channel: %s\n", channel)
}

func (c *ChatClient) ReadMessages() {
	msgs, err := c.ch.Consume(
		c.currentQueue,
		"",    // consumer
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		log.Printf("Consume error: %v", err)
		return
	}

	for msg := range msgs {
		fmt.Printf("\n[%s] %s\n> ", msg.RoutingKey, msg.Body)
	}
}

func (c *ChatClient) HandleInput() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "!switch ") {
			parts := strings.SplitN(text, " ", 2)
			if len(parts) == 2 {
				c.SwitchChannel(parts[1])
			}
		} else {
			err := c.ch.Publish(
				exchangeName,
				c.currentChan,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(text),
				},
			)
			if err != nil {
				log.Printf("Publish error: %v", err)
			}
		}
		fmt.Print("> ")
	}
}

func (c *ChatClient) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
