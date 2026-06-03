package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/games"
)

func main() {
	results, err := games.ItchSearch("car", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d game(s)\n\n", len(results))
	for i, g := range results {
		rating := "no rating"
		if g.Rating != nil {
			rating = fmt.Sprintf("%.2f (%d ratings)", g.Rating.Average, g.Rating.Total)
		}
		fmt.Printf("[%d] %s by %s\n    %s\n    Rating: %s | Genre: %s\n\n",
			i+1, g.Title, g.Author, g.URL, rating, g.Genre)
	}
}
