package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/scrapers"
)

func main() {
	results, err := scrapers.DuckDuckGoSearch("golang", 2)
	if err != nil {
		log.Fatal(err)
	}
	for i, r := range results {
		fmt.Printf("[%d] %s\n    %s\n    %s\n\n", i+1, r.Title, r.URL, r.Snippet)
	}
}
