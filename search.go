package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

const MAX_RESULTS = 5
const MIN_RESULTS = 1

func runQuery(es *elasticsearch.Client, searchIndex string, searchTerms string) (string, error) {

	trimmed := strings.TrimSpace(searchTerms)
	query := fmt.Sprintf(`{
  "retriever": {
    "rrf": {
      "retrievers": [
        {
          "standard": {
            "query": {
              "semantic": {
                "field": "semantic_title",
                "query": "%s"
              }
            }
          }
        },
        {
          "standard": {
            "query": {
              "semantic": {
                "field": "semantic_wysd",
                "query": "%s"
              }
            }
          }
        },
        {
          "standard": {
            "query": {
              "multi_match": {
                "query": "{query}",
                "fields": [
                  "Product",
                  "Title",
                  "Issue",
                  "Organization",
                  "Category"
                ]
              }
            }
          }
        }
      ],
      "rank_constant": 20,
				"rank_window_size": 50
    }
  },
  "highlight": {
    "fields": {
      "semantic_title": {
        "type": "semantic",
        "number_of_fragments": 2,
        "order": "score"
      },
      "semantic_wysd": {
        "type": "semantic",
        "number_of_fragments": 3,
        "order": "score"
      }
    }
  }
}`, trimmed, trimmed)

	returned, err := es.Search(
		es.Search.WithIndex(searchIndex),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithFrom(MIN_RESULTS),
		es.Search.WithSize(MAX_RESULTS),
		es.Search.WithPretty(),
	)
	if err != nil {
		return "", err
	}

	if returned.Status() != "200 OK" {
		fmt.Printf("\n\nQUERY ERROR: \n\n%s", returned)
		log.Fatal(returned)
	}

	defer returned.Body.Close()
	result, err := io.ReadAll(returned.Body)

	return string(result), err
}
