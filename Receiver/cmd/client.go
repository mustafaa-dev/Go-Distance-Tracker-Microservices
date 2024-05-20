package main

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)

	// Subscribe to subject
	sub, _ := nc.SubscribeSync("otu")

	// Use a Go routine to keep reading messages until the program is stopped
	go func() {
		for {
			// Create a context with no deadline
			ctx := context.Background()

			// Wait for the next message
			msg, err := sub.NextMsgWithContext(ctx)
			if err != nil {
				log.Println("Error reading message:", err)
				continue
			}

			log.Println("Received message:", string(msg.Data))
		}
	}()

	// Keep the connection alive until the program is stopped
	select {}
}
