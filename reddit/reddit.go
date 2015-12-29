package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type RedditComment struct {
	Data struct {
		Children  []*RedditComment `json:"children"`
		Author    string           `json:"author"`
		Permalink string           `json:"permalink"`
		Title     string           `json:"title"`
		URL       string           `json:"url"`
	} `json:"data"`
}

type RedditIndexPage struct {
	Data struct {
		Children []*RedditComment `json:"children"`
	} `json:"data"`
}

type RedditArticlePage []struct {
	Data struct {
		Children []*RedditComment `json:"children"`
	} `json:"data"`
}

type JsonPermalink string

// Retrieve the main HOT listings page
func GetGolangIndex() []*JsonPermalink {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://www.reddit.com/r/golang/hot.json", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "/u/jbuberel")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Unable to retrieve /r/golang index")
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	bodyText := string(body[:])
	log.Println(len(bodyText))

	var listings RedditIndexPage
	err = json.Unmarshal(body, &listings)
	if err != nil {
		log.Println("error:", err)
	}

	var permalinks = make([]*JsonPermalink, 0)
	for i, listing := range listings.Data.Children {
		log.Printf("%v - %v\n", i, listing.Data.Title)
		log.Printf("   %v\n", listing.Data.Permalink)
		var permalink JsonPermalink
		permalink = JsonPermalink("http://www.reddit.com" + listing.Data.Permalink[:len(listing.Data.Permalink)-1] + ".json")
		log.Printf("  %v\n", permalink)
		permalinks = append(permalinks, &permalink)
	}
	return permalinks

}

// Author The name of the person who wrote the comment
type Author string

// GetRedditAuthors Get the list of authors
func GetRedditAuthors(permalink JsonPermalink) []Author {

	log.Printf("  permalink: %v\n", permalink)

	client := &http.Client{}

	req, err := http.NewRequest("GET", string(permalink), nil)
	if err != nil {
		log.Printf("Unable to retrieve %v\n", permalink)
		return make([]Author, 0)
	}

	req.Header.Set("User-Agent", "/u/jbuberel")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Unable to retrieve /r/golang index")
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	bodyText := string(body[:])
	log.Println(len(bodyText))

	var authors = make([]Author, 0)

	re1, err := regexp.Compile(`\"author\":\s*\"([^\"]+)\"`)
	result := re1.FindAllStringSubmatch(bodyText, -1)

	for k, v := range result {
		fmt.Printf("%d. %s\n", k, v[1])
		authors = append(authors, Author(v[1]))
	}

	return authors

}

// Capture Start the capture
func Capture(dirname *string) {
	log.Println("starting reddit capture")
	permalinks := GetGolangIndex()

	var authorMap = make(map[Author]int)
	for _, permalink := range permalinks {
		fmt.Println(*permalink)
		authors := GetRedditAuthors(*permalink)

		for _, author := range authors {
			authorMap[author]++
		}
	}

	timestamp := time.Now().Format("2006-01-02")
	outfile := fmt.Sprintf("%v/reddit-authors-%v.csv", *dirname, timestamp)
	fmt.Printf("Saving results to file %v\n", outfile)
	f, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("Unable to create file %v - %v\n", outfile, err)
	}
	defer f.Close()

	for k, v := range authorMap {
		fmt.Printf("%v - %v\n", k, v)
		f.WriteString(fmt.Sprintf("%v,%v,%v\n", timestamp, k, v))

	}
	fmt.Printf("Unique authors: %v\n", len(authorMap))

}
