package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/go-tdx-guest/client"
)

func main() {
	var quoteProvider interface{}
	var err error

	// Try to get the quote provider first
	quoteProvider, err = client.GetQuoteProvider()
	if err != nil {
		log.Printf("Failed to get quote provider: %v", err)
		log.Println("Falling back to TDX Guest device...")

		// Fall back to opening the TDX Guest device
		quoteProvider, err = client.OpenDevice()
		if err != nil {
			log.Fatalf("Failed to open TDX Guest device: %v", err)
		}
		if closer, ok := quoteProvider.(client.Device); ok {
			defer closer.Close()
		}
	}

	// Prepare the report data with "Hello from TDX VM"
	message := "Hello from TDX VM"
	var reportData [64]byte
	copy(reportData[:], message)

	// Get the raw quote
	rawQuote, err := client.GetRawQuote(quoteProvider, reportData)
	if err != nil {
		log.Fatalf("Failed to get raw quote: %v", err)
	}

	// Write raw quote to file
	err = os.WriteFile("quote.dat", rawQuote, 0644)
	if err != nil {
		log.Fatalf("Failed to write raw quote to file: %v", err)
	}

	// Get the quote in proto format
	quote, err := client.GetQuote(quoteProvider, reportData)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
	}

	// Marshal the proto quote to JSON
	quoteJSON, err := json.MarshalIndent(quote, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal quote to JSON: %v", err)
	}

	// Print the quote information
	fmt.Printf("Successfully generated quote:\n")
	fmt.Printf("Report Data: %s\n", string(reportData[:len(message)]))
	fmt.Printf("Report Data (hex): %s\n", hex.EncodeToString(reportData[:]))
	fmt.Printf("Raw Quote (first 32 bytes): %s\n", hex.EncodeToString(rawQuote[:32]))
	fmt.Printf("Raw Quote length: %d bytes\n", len(rawQuote))
	fmt.Printf("Raw Quote written to quote.dat\n")
	fmt.Printf("Quote (JSON):\n%s\n", string(quoteJSON))

	// Determine which method was used
	switch quoteProvider.(type) {
	case client.QuoteProvider:
		fmt.Println("Quote obtained using QuoteProvider")
	case client.Device:
		fmt.Println("Quote obtained using Device (fallback method)")
	default:
		fmt.Println("Unknown quote provider type")
	}
}
