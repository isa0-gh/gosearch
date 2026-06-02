package main

import (
	"fmt"
	"log"
	"time"

	"github.com/isa0-gh/gosearch/internal/apps"
)

func main() {
	results, err := apps.FlathubSearch("browser", 1)
	if err != nil {
		log.Fatal(err)
	}
	for i, a := range results {
		fmt.Printf("[%d] %s (%s)\n    %s\n    Developer: %s  License: %s\n    Updated: %s\n    %s\n\n",
			i+1, a.Name, a.AppID, a.Summary, a.Developer, a.License,
			time.Unix(a.UpdatedAt, 0).Format("2006-01-02"), a.URL)
	}
}
