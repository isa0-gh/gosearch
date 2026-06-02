package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/torrents"
)

func main() {
	results, err := torrents.NyaaSearch("naruto", 2)
	if err != nil {
		log.Fatal(err)
	}
	for i, t := range results {
		fmt.Printf("[%d] %s\n    Category: %s  Size: %s  Seeds: %d  Leeches: %d\n    %s\n\n",
			i+1, t.Name, t.Category, t.Size, t.Seeders, t.Leechers, t.URL)
	}
}
