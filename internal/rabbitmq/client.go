package rabbitmq

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	exchangeName = "chat_exchange"
)

type ChatClient struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	currentChan string
	consumerTag string

	Messages *widget.Entry
	Type     string

	clientID string

	consumerCancel context.CancelFunc
	mutex          sync.Mutex
}

func (c *ChatClient) Connect(addr string) error {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://admin:admin@%s:5672/", addr))
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
		true,  // durable
		false, // autoDelete
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

func NewClient() *ChatClient {
	return &ChatClient{
		clientID: uuid.New().String(),
	}
}

func (c *ChatClient) SwitchChannel(channel string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.consumerCancel != nil {
		c.consumerCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.consumerCancel = cancel

	// Создаем уникальную очередь для каждого клиента
	queue, err := c.ch.QueueDeclare(
		"",    // имя генерируется автоматически
		false, // durable
		true,  // autoDelete
		true,  // exclusive
		false,
		nil,
	)
	if err != nil {
		log.Printf("Queue declare error: %v", err)
		return
	}

	// Привязываем очередь к каналу
	err = c.ch.QueueBind(
		queue.Name,
		channel+".#", // слушаем все сообщения в канале
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Queue bind error: %v", err)
		return
	}

	if c.Type == "gui" {
		go c.ReadMessagesTo(ctx, queue.Name)
	} else {
		go c.ReadMessages(ctx, queue.Name)
	}
	c.currentChan = channel
	fmt.Printf("Switched to channel: %s\n", channel)
}

func (c *ChatClient) ReadMessages(ctx context.Context, queueName string) {
	msgs, err := c.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Consume error: %v", err)
		return
	}

	for {
		select {
		case msg := <-msgs:
			senderID, ok := msg.Headers["sender_id"].(string)
			if ok && senderID == c.clientID {
				continue
			}
			fmt.Printf("\n[%s] %s\n> ", msg.RoutingKey, msg.Body)
		case <-ctx.Done():
			return
		}
	}
}

func (c *ChatClient) ReadMessagesTo(ctx context.Context, queueName string) {
	msgs, err := c.ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Consume error: %v", err)
		return
	}

	for {
		select {
		case msg := <-msgs:
			c.Messages.SetText(c.Messages.Text + fmt.Sprintf("[%s] %s\n", msg.RoutingKey, msg.Body))
		case <-ctx.Done():
			return
		}
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
			err := c.SendMessage(text)
			if err != nil {
				log.Printf("Publish error: %v", err)
			}
		}
		fmt.Print("> ")
	}
}

func (c *ChatClient) SendMessage(message string) error {
	fullRoutingKey := c.currentChan + ".message"
	return c.ch.Publish(
		exchangeName,   // exchange
		fullRoutingKey, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (c *ChatClient) Close() {
	if c.consumerCancel != nil {
		c.consumerCancel()
	}
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
