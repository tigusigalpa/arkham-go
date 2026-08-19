package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// Fetch enriched transfers from Binance
	filter := &arkham.TransferFilter{
		Base:  []string{"binance"},
		Limit: 10,
	}

	transfers, meta, err := client.Transfers.Transfers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Transfers failed: %v", err)
	}

	for i, tx := range transfers.Transfers {
		fmt.Printf("%d: %s -> %s  $%.2f  (%s)\n", i+1, tx.From, tx.To, tx.USD, tx.Time)
	}
	fmt.Printf("\nTotal: %d transfers\n", len(transfers.Transfers))
	fmt.Printf("Intel remaining: %s\n", meta.IntelDatapointsRemaining)
}
