// Track daily list of importers of the golang.org/x/mobile/app package:
//   http://godoc.org/golang.org/x/mobile/app?importers
package mobile

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var searchUrl = "http://godoc.org/golang.org/x/mobile/app?importers"

func Capture(dirname *string) {
	doc, err := goquery.NewDocument(searchUrl)
	if err != nil {
		log.Fatal(err)
	}

	matches := 0
	doc.Find("div.container table tr>td").Each(func(i int, tr *goquery.Selection) {
		htmlText, err := tr.First().Html()
		if err == nil && strings.Contains(htmlText, "/") && !strings.Contains(htmlText, "golang.org/x/mobile") {
			fmt.Printf("Review %d: %v\n", i, htmlText)
			matches++
		}
	})

	timestamp := time.Now().Format("2006-01-02")
	outfile := fmt.Sprintf("%v/mobile-imports-%v.csv", *dirname, timestamp)
	log.Printf("Saving results to file %v\n", outfile)
	f, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("Unable to create file %v - %v\n", outfile, err)
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%v,%v\n", timestamp, matches))

}
