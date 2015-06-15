package reddit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"fmt"
	"os"
)

type RedditComment struct {
	Kind    string `json:"kind"`
	After   string `json:"after"`
	Before  string `json:"before"`
	Modhash string `json:"modhash"`
	Data    struct {
		Children     []*RedditComment `json:"children"`
		ApprovedBy   string           `json:"approved_by"`
		Author       string           `json:"author"`
		BannedBy     string           `json:"banned_by"`
		Clicked      bool             `json:"clicked"`
		Created      float64          `json:"created"`
		CreatedUtc   float64          `json:"created_utc"`
		Domain       string           `json:"domain"`
		Downs        int              `json:"downs"`
		Gilded       int              `json:"gilded"`
		ID           string           `json:"id"`
		Name         string           `json:"name"`
		NumComments  int              `json:"num_comments"`
		Permalink    string           `json:"permalink"`
		Score        int              `json:"score"`
		Selftext     string           `json:"selftext"`
		SelftextHTML string           `json:"selftext_html"`
		Subreddit    string           `json:"subreddit"`
		SubredditID  string           `json:"subreddit_id"`
		Thumbnail    string           `json:"thumbnail"`
		Title        string           `json:"title"`
		Ups          int              `json:"ups"`
		URL          string           `json:"url"`
	} `json:"data"`
}

type RedditIndexPage struct {
	Data struct {
		Kind     string           `json:"kind"`
		After    string           `json:"after"`
		Before   string           `json:"before"`
		Modhash  string           `json:"modhash"`
		Children []*RedditComment `json:"children"`
	} `json:"data"`
}

type RedditArticlePage []struct {
	Data struct {
		Kind     string           `json:"kind"`
		After    string           `json:"after"`
		Before   string           `json:"before"`
		Modhash  string           `json:"modhash"`
		Children []*RedditComment `json:"children"`
	} `json:"data"`
}

type JsonPermalink string

func GetGolangIndex() []*JsonPermalink {

	resp, err := http.Get("http://www.reddit.com/r/golang/hot.json")
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

type RedditAuthor string

func GetRedditAuthors(permalink JsonPermalink) []*RedditAuthor {
	var authors = make([]*RedditAuthor, 0)

	resp, err := http.Get(string(permalink))
	if err != nil {
		log.Printf("Unable to retrieve %v\n", permalink)
		return authors
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	bodyText := string(body[:])
	log.Println(len(bodyText))

	var articlePage RedditArticlePage
	err = json.Unmarshal(body, &articlePage)
	if err != nil {
		log.Println("error:", err)
	}
	for _, article := range articlePage {
		for _, listing := range article.Data.Children {
			var author RedditAuthor = RedditAuthor(listing.Data.Author)
			authors = append(authors, &author)
			if len(listing.Data.Children) > 0 {
				sublink := JsonPermalink("http://www.reddit.com" + listing.Data.Permalink[:len(listing.Data.Permalink)-1] + ".json")
				authors = append(authors, GetRedditAuthors(sublink)...)
			}
		}
	}

	return authors

}

func Capture(dirname *string) {
	log.Println("starting reddit capture")
	permalinks := GetGolangIndex()

	var authorMap = make(map[RedditAuthor]int)
	for _, permalink := range permalinks {
		fmt.Println(*permalink)
		authors := GetRedditAuthors(*permalink)

		for _, author := range authors {
			authorMap[*author] += 1
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
