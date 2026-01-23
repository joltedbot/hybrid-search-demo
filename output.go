package main

import "fmt"

func printResults(results Result) {
	resultCount := results.Hits.Total.Value
	hitCount := len(results.Hits.Hits)

	for i := range hitCount - 1 {
		source := results.Hits.Hits[i].Source
		score := results.Hits.Hits[i].Score
		fmt.Printf("Score:     %f\n", score)
		fmt.Printf("Product:   %s\n", source["Product"])
		fmt.Printf("Title:     %s\n", source["Title"])
		fmt.Printf("URL:       %s\n", source["URL"])
		fmt.Printf("What you should do:\n%s\n", source["What you should do"])
		fmt.Printf("\n-------------\n")
	}

	fmt.Printf("\n\nReturned top %d of the %d results\n\n", hitCount, resultCount)
}
