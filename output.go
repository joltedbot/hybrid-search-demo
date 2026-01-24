package main

import "fmt"

func printResults(results Result) {
	hitCount := len(results.Hits.Hits)

	for i := range hitCount - 1 {
		source := results.Hits.Hits[i].Source

		fmt.Printf("Title:     %s\n", source["Title"])
		if source["Product"] != nil {
			fmt.Printf("Product:   %s\n", source["Product"])
		} else {
			fmt.Println("Product:   N/A")
		}
		fmt.Printf("URL:       %s\n", source["URL"])
		fmt.Printf("What you should do:\n%s\n", source["What you should do"])
		fmt.Printf("\n-------------\n")
	}

	fmt.Printf("\nReturned %d Results\n\n", hitCount)

}
