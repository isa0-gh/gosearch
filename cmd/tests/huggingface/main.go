package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/ml"
)

func main() {
	fmt.Println("Searching Hugging Face for 'minimax'...")
	results, err := ml.HuggingFaceSearch("minimax", 1)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	fmt.Printf("Found %d results:\n\n", len(results))
	for i, r := range results {
		fmt.Printf("[%d] Name: %s\n", i+1, r.Name)
		fmt.Printf("    URL:  %s\n", r.URL)
		fmt.Printf("    Desc: %s\n", r.Description)
		fmt.Printf("    Caps: %v\n", r.Capabilities)
		fmt.Printf("    Size: %s\n", r.Size)
		fmt.Printf("    Pulls: %s\n", r.Pulls)
		fmt.Printf("    Tags: %s\n", r.Tags)
		fmt.Printf("    Updated: %s\n\n", r.Updated)
		if i >= 4 {
			break
		}
	}
}
