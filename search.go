package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

const SEMENTIC_TEXT_RANK_BOOST = 0.75
const MININUM_SCORE = 0.5

func runQuery(es *elasticsearch.Client, searchIndex string, searchTerms string) (string, error) {

	trimmed := strings.TrimSpace(searchTerms)
	query := fmt.Sprintf(`{
		"retriever": {
			"rrf": {
				"retrievers": [
					{
						"standard": {
							"query": {
								"match": {
									"Product": "%s"
								}
							}
						}
					},
					{
						"standard": {
							"min_score": %f,
							"query": {
								"semantic": {
									"field": "semantic_text",
									"query": "%s",
									"boost": %f
								}
							}
						}
					}
				]
			}
		}
	}`, trimmed, MININUM_SCORE, trimmed, SEMENTIC_TEXT_RANK_BOOST)

	returned, err := es.Search(
		es.Search.WithIndex(searchIndex),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithTrackTotalHits(true),
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
