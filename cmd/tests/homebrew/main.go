package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/apps"
)

func main() {
	results, err := apps.HomebrewSearch("wget", 1)
	if err != nil {
		log.Fatal(err)
	}
	for i, a := range results {
		fmt.Printf("[%d] %s (%s)\n    %s\n    Type: %s  License: %s\n    %s\n\n",
			i+1, a.Name, a.AppID, a.Summary, a.Developer, a.License, a.URL)
	}

	fmt.Println("--- CASK TEST ---")
	results, err = apps.HomebrewSearch("firefox", 1)
	if err != nil {
		log.Fatal(err)
	}
	for i, a := range results {
		fmt.Printf("[%d] %s (%s)\n    %s\n    Type: %s  License: %s\n    %s\n\n",
			i+1, a.Name, a.AppID, a.Summary, a.Developer, a.License, a.URL)
	}
}
