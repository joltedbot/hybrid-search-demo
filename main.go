package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	apiKey := os.Getenv("ES_API_KEY")
	esUrl := os.Getenv("ES_SERVER_URL")
	searchIndex := os.Getenv("ES_SEARCH_INDEX")

	if apiKey == "" || esUrl == "" || searchIndex == "" {
		fmt.Println("Please set environmental variables ES_API_KEY, ES_SERVER_URL, ES_SEARCH_INDEX and try again.")
		os.Exit(1)
	}

	esClient, err := setupElasticsearch(esUrl, apiKey)
	if err != nil {
		log.Fatalf("Error setting up Elasticsearch: %s", err)
	}

	m := initialModel(esClient, searchIndex)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}

func setupElasticsearch(url, apiKey string) (*elasticsearch.Client, error) {
	config := elasticsearch.Config{
		Addresses: []string{url},
		APIKey:    apiKey,
	}
	return elasticsearch.NewClient(config)
}
