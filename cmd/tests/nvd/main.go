package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/vuln"
)

func main() {
	results, err := vuln.NVDSearch("python", 1)
	if err != nil {
		log.Fatal(err)
	}
	for i, c := range results {
		fmt.Printf("[%d] %s  [%s %.1f]\n    %s\n    %s\n\n",
			i+1, c.ID, c.Severity, c.Score, c.Description, c.URL)
	}
}
