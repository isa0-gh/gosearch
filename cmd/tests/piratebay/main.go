package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/torrents"
)

func main() {
	results, err := torrents.PirateBaySearch("undertale", 0)
	if err != nil {
		log.Fatal(err)
	}
	for i, t := range results {
		fmt.Printf("[%d] %s\n    Seeds: %d  Leeches: %d  Size: %d bytes\n    Category: %s  Uploader: %s\n    %s\n\n",
			i+1, t.Name, t.Seeders, t.Leechers, t.Size, t.Category, t.Uploader, t.MagnetURL[:60]+"...")
	}
}
