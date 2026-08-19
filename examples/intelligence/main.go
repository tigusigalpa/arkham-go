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

	// Look up intelligence for an address
	addr, meta, err := client.Intelligence.Address(context.Background(), "0x28C6c06298d514Db089934071355E5743bf21d60", nil)
	if err != nil {
		log.Fatalf("Intelligence.Address failed: %v", err)
	}

	fmt.Printf("Address: %s\n", addr.Address)
	fmt.Printf("Chain:   %s\n", addr.Chain)
	if addr.ArkhamEntity != nil {
		fmt.Printf("Entity:  %s\n", addr.ArkhamEntity.Name)
	}
	fmt.Printf("Intel usage: %s / %s\n", meta.IntelDatapointsUsage, meta.IntelDatapointsLimit)
}
