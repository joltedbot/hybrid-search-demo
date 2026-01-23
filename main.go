package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

func main() {

	apiKey := os.Getenv("ES_API_KEY")
	esUrl := os.Getenv("ES_SERVER_URL")
	searchIndex := os.Getenv("ES_SEACH_INDEX")

	if apiKey == "" || esUrl == "" {
		fmt.Println("Please set environmental variables ES_API_KEY, ES_SERVER_URL, ES_SEARCH_INDEX to your ES API key, ES Server's URL, and index to search and try again")
		os.Exit(1)
	}

	configuration := elasticsearch.Config{
		Addresses: []string{
			esUrl,
		},
		APIKey: apiKey,
	}

	es, err := elasticsearch.NewClient(configuration)
	if err != nil {
		log.Fatal(err)
	}

	defer es.Close(context.Background())

	/*
		connectionInfo, err := es.Info()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nSuccessfully Connected:\n%s\n\n", connectionInfo)
	*/

	fmt.Printf("\nSuccessfully Connected\n\n")

	fmt.Println("\n\nWelcome to the Hybrid Search Demo Tool!")

	for {
		fmt.Println("\n What product would you like to search for?:")
		fmt.Print("> ")
		buffer := bufio.NewReader(os.Stdin)
		searchTerms, err := buffer.ReadString('\n')
		if err != nil {
			fmt.Println("Bad search terms. Try again!")
			continue
		}

		fmt.Printf("Retrieving results for: %s\n", strings.TrimSpace(searchTerms))

		queryResult, err := runQuery(es, searchIndex, searchTerms)
		if err != nil {
			log.Fatal(err)
		}

		results := Result{}
		err = json.Unmarshal([]byte(queryResult), &results)
		if err != nil {
			log.Fatal(err)
		}

		printResults(results)

	}

}
