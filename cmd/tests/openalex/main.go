package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/academic"
)

func main() {
	results, err := academic.OpenAlexSearch("python", 1)
	if err != nil {
		log.Fatal(err)
	}
	for i, p := range results {
		fmt.Printf("[%d] %s\n    Authors: %s\n    Type: %s\n    %s\n    %s\n\n",
			i+1, p.Title, p.Authors, p.Type, p.Abstract, p.URL)
	}
}
