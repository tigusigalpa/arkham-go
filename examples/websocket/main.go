package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tigusigalpa/arkham-go"
)

func main() {
	apiKey := os.Getenv("ARKHAM_API_KEY")
	if apiKey == "" {
		log.Fatal("ARKHAM_API_KEY not set")
	}

	client, err := arkham.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Step 1: Create a WebSocket v2 stream
	req := &arkham.CreateStreamV2Request{
		Base:   []string{"binance"},
		UsdGte: "500000",
	}

	stream, _, err := client.Streams.Create(context.Background(), req)
	if err != nil {
		log.Fatalf("Create stream failed: %v", err)
	}
	fmt.Printf("Stream created: %s\n", stream.StreamID)

	// Ensure cleanup
	defer func() {
		_, _, err := client.Streams.Delete(context.Background(), stream.StreamID)
		if err != nil {
			log.Printf("Delete stream failed: %v", err)
		}
	}()

	// Step 2: Connect to the stream
	conn, err := client.Streams.Connect(context.Background(), stream.StreamID)
	if err != nil {
		log.Fatalf("Connect failed: %v", err)
	}
	defer conn.Close()

	// Step 3: Receive transfers for 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	for {
		msg, err := conn.ReceiveTyped()
		if err != nil {
			log.Printf("Receive ended: %v", err)
			return
		}
		fmt.Printf("[%s] %s\n", msg.Type, string(msg.Payload))
	}
}
