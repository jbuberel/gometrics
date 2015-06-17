package github

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
	TotalCount int           `json:"total_count"`
	Repos      []*GithubRepo `json:"items"`
}

type GithubRepo struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	FullName        string  `json:"full_name"`
	Url             string  `json:"url"`
	Size            int     `json:"size"`
	StargazerCount  int     `json:"stargazers_count"`
	WatchersCount   int     `json:"watchers_count"`
	Language        string  `json:"language"`
	ForksCount      int     `json:"forks_count"`
	OpenIssuesCount int     `json:"open_issues_count"`
	Score           float64 `json:"score"`
}

var searchUrl = "https://api.github.com/search/repositories"

func getResults(sortyby	string) []*GithubSearchResult {
	searchBase := searchUrl + fmt.Sprintf("?q=language:go&sort=%v&order=desc", sortyby)
	var results []*GithubSearchResult = make([]*GithubSearchResult, 0)
	for page := 1; page < 10; page++ {
		search := searchBase + fmt.Sprintf("&page=%v", page)
		log.Printf("Searching for URL: %v\n", search)

		time.Sleep(5 * time.Second)
		resp, err := http.Get(string(search))
		if err != nil {
			log.Printf("Unable to retrieve %v\n", search)
			return results
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		bodyText := string(body[:])
		log.Println(len(bodyText))

		var result GithubSearchResult
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Println("error:", err)
		}
		results = append(results, &result)
		for _, repo := range result.Repos {
			log.Printf("  repo: %v - %v\n", repo.Id, repo.FullName)
		}

	}

	return results
}

func getStarred() []*GithubSearchResult {
	return getResults("stars")	
}

func getForked() []*GithubSearchResult  {
	return getResults("forks")	

}

func Capture(dirname *string, githubClientId, githubSecretKey string, githubSecretToken string) {
	timestamp := time.Now().Format("2006-01-02")
	searchResults := getStarred()
	outfile := fmt.Sprintf("%v/github-starred-%v.csv", *dirname, timestamp)
	log.Printf("Saving results to file %v\n", outfile)
	f, err := os.Create(outfile)
	if err != nil {
		fmt.Printf("Unable to create file %v - %v\n", outfile, err)
	}
	defer f.Close()

	for _, sr := range searchResults {
		for _, repo := range sr.Repos {
			log.Printf("%v - %v\n", repo.Id, repo.Name)
			f.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", timestamp, sr.TotalCount, repo.Id, repo.Name, repo.FullName,
				repo.Url, repo.Size, repo.StargazerCount, repo.WatchersCount, repo.Language, repo.ForksCount, repo.OpenIssuesCount,
				repo.Score))
		}
	}
	
	searchResults = getForked()
	outfile = fmt.Sprintf("%v/github-forked-%v.csv", *dirname, timestamp)
	log.Printf("Saving results to file %v\n", outfile)
	f, err = os.Create(outfile)
	if err != nil {
		fmt.Printf("Unable to create file %v - %v\n", outfile, err)
	}
	defer f.Close()

	for _, sr := range searchResults {
		for _, repo := range sr.Repos {
			log.Printf("%v - %v\n", repo.Id, repo.Name)
			f.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", timestamp, sr.TotalCount, repo.Id, repo.Name, repo.FullName,
				repo.Url, repo.Size, repo.StargazerCount, repo.WatchersCount, repo.Language, repo.ForksCount, repo.OpenIssuesCount,
				repo.Score))
		}
	}


}
