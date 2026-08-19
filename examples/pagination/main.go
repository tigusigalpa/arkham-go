package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

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

	// Build query from filter
	q := url.Values{}
	q.Set("from", "0x28C6c06298d514Db089934071355E5743bf21d60")
	q.Set("limit", strconv.Itoa(100))

	paginator := arkham.NewPaginator(
		context.Background(),
		client,
		"/transfers",
		q,
		100,  // page size
		1000, // max items
		10,   // max pages
	)

	total := 0
	for paginator.HasNext() {
		var resp arkham.EnrichedTransfers
		_, err := paginator.NextPage(&resp)
		if err != nil {
			log.Fatalf("Page fetch failed: %v", err)
		}
		if len(resp.Transfers) == 0 {
			break
		}
		total += len(resp.Transfers)
		fmt.Printf("Got %d transfers (total: %d)\n", len(resp.Transfers), total)
	}

	fmt.Printf("\nDone. Total transfers: %d\n", total)
}
