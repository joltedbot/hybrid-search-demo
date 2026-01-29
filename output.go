package main

import (
	"fmt"
	"strings"
)

func printResults(results Result) {
	hitCount := len(results.Hits.Hits)

	for i := range hitCount {
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

func showUIResults(results Result) string {
	hitCount := len(results.Hits.Hits)
	var builder strings.Builder

	for i := range hitCount {
		source := results.Hits.Hits[i].Source

		builder.WriteString(fmt.Sprintf("Title:     %s\n", source["Title"]))
		if source["Product"] != nil {
			builder.WriteString(fmt.Sprintf("Product:   %s\n", source["Product"]))
		} else {
			builder.WriteString("Product:   N/A\n")
		}
		builder.WriteString(fmt.Sprintf("URL:       %s\n", source["URL"]))
		builder.WriteString(fmt.Sprintf("What you should do:\n%s\n", source["What you should do"]))
		builder.WriteString("\n-------------\n")
	}
	builder.WriteString(fmt.Sprintf("\nReturned %d Results\n\n", hitCount))

	return builder.String()
}
