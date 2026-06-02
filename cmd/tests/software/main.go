package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/software"
)

func main() {
	gh, err := software.GitHubSearch("golang web framework", 1)
	if err != nil {
		log.Fatal("github:", err)
	}
	fmt.Printf("GitHub: %d results\n", len(gh))
	if len(gh) > 0 {
		fmt.Printf("  [1] %s - %s (%d★)\n", gh[0].Name, gh[0].URL, gh[0].Stars)
	}

	gl, err := software.GitLabSearch("golang", 1)
	if err != nil {
		log.Fatal("gitlab:", err)
	}
	fmt.Printf("GitLab: %d results\n", len(gl))
	if len(gl) > 0 {
		fmt.Printf("  [1] %s - %s (%d★)\n", gl[0].Name, gl[0].URL, gl[0].Stars)
	}

	sf, err := software.SourceForgeSearch("golang", 1)
	if err != nil {
		log.Fatal("sourceforge:", err)
	}
	fmt.Printf("SourceForge: %d results\n", len(sf))
	if len(sf) > 0 {
		fmt.Printf("  [1] %s - %s\n", sf[0].Name, sf[0].URL)
	}
}
