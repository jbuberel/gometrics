package gitissues

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type GithubSearchResult struct {
	Id        int    `json:"id"`
	Number    int    `json:"number"`
	State     string `json:"state"`
	Title     string `json:"title"`
	ClosedAt  string `json:"closed_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Milestone struct {
		Number int    `json:"number"`
		Title  string `jason:"number"`
	} `json:"milestone"`
}

type GithubRepo struct {
}

var searchUrl = "https://api.github.com/repos/golang/go/issues?state=open"

func getResults(githubSecretToken string) []GithubSearchResult {
	var results []GithubSearchResult = make([]GithubSearchResult, 0)
	found := false
//        for page := 1; page < 5; page++ {
	for page := 1; page == 1 || found; page++ {
		search := searchUrl + fmt.Sprintf("&page=%v", page)
		log.Printf("Searching for URL: %v\n", search)
		found = false
		time.Sleep(1 * time.Second)

		client := &http.Client{}

		req, err := http.NewRequest("GET", string(search), nil)
		if err != nil {
			log.Println("Unable to create HTTP request", err)
			return results
		}
		req.SetBasicAuth("jbuberel", githubSecretToken)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Unable to retrieve %v\n", search)
			return results
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		bodyText := string(body[:])
		log.Println(len(bodyText))

		var result []GithubSearchResult
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("error unmarshalling %v of data [%s]\n", err, body)
		}
		results = append(results, result...)
		for _, issue := range result {
			log.Printf("  issue: %v - %v - %v\n", issue.Id, issue.Milestone.Title, issue.Title)
			found = true
		}

	}

	return results
}

func Capture(dirname *string, githubClientId, githubSecretKey string, githubSecretToken string) {
	timestamp := time.Now().Format("2006-01-02")
	searchResults := getResults(githubSecretToken)
	outfile := fmt.Sprintf("%v/github-issues-%v.csv", *dirname, timestamp)
	log.Printf("Saving results to file %v\n", outfile)
	f, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("Unable to create file %v - %v\n", outfile, err)
	}
	defer f.Close()

	total := 0
	issuesByMilestone := make(map[string]int)
	for _, issue := range searchResults {
		log.Printf("%v - %v - %v\n", issue.Id, issue.Milestone.Title, issue.Title)
		total += 1
		issuesByMilestone[issue.Milestone.Title] += 1
	}

	for m := range issuesByMilestone {
		f.WriteString(fmt.Sprintf("%v,%v,%v,%v\n", timestamp, total, m, issuesByMilestone[m]))
	}
}
