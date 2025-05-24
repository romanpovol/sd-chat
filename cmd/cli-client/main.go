package main

import (
	"flag"
	"log"

	"github.com/romanpovol/sd-chat/internal/rabbitmq"
)

func main() {
	serverAddr := flag.String("server", "127.0.0.1", "RabbitMQ Server address")
	initChannel := flag.String("channel", "general", "Initial channel")
	flag.Parse()

	client := rabbitmq.NewClient()
	defer client.Close()

	err := client.Connect(*serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client.SwitchChannel(*initChannel)

	client.HandleInput()
}
