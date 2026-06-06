package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/isa0-gh/gosearch/internal/games"
)

func main() {
	results, err := games.GogSearch("witcher", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d game(s)\n\n", len(results))
	for i, g := range results {
		fmt.Printf("[%d] %s\n    %s\n    %s | Platforms: %s | Genres: %s\n\n",
			i+1, g.Title, g.URL, g.Price, strings.Join(g.Platforms, ", "), strings.Join(g.Genres, ", "))
	}
}
