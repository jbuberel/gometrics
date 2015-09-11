package twittermobile

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jbuberel/anaconda"
)

func extract(api *anaconda.TwitterApi, term string, since *string, until *string) map[string]anaconda.Tweet {
	log.Printf("Beginning #golang #mobile extraction for date range %v to %v.\n", *since, *until)
	v := url.Values{}
	v.Add("result_type", "recent")

	tweets := make(map[string]anaconda.Tweet)

	q := fmt.Sprintf("%v since:%v until:%v", term, *since, *until)
	log.Printf("q: %v\n", q)
	for searchResult, _ := api.GetSearch(q, v); len(searchResult.Statuses) > 0; searchResult, _ = searchResult.GetNext(api) {
		for _, tweet := range searchResult.Statuses {
			tweets[tweet.IdStr] = tweet
		}
		log.Printf("Total tweets %v\n", len(tweets))
		time.Sleep(5 * time.Second)
	}

	log.Printf("Completing tweet extraction, found %v tags and %v mentions.", len(tweets), 0)
	return tweets

}

// Capture captures the twitter data for #golang #mobile
func Capture(dirname *string, since *string, until *string, twitterConsumerKey, twitterConsumerSecret, twitterAccessToken, twitterSecretToken string) {
	log.Printf("Connecting to twitter\n")
	anaconda.SetConsumerKey(twitterConsumerKey)
	anaconda.SetConsumerSecret(twitterConsumerSecret)
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterSecretToken)
	timestamp := time.Now().Format("2006-01-02")

	tweets := extract(api, "#golang #mobile", since, until)
	outfile := fmt.Sprintf("%v/twitter-mobile-%v.csv", *dirname, timestamp)
	log.Printf("Saving results to file %v\n", outfile)
	f, err := os.Create(outfile)
	if err != nil {
		log.Printf("Unable to create file %v - %v\n", outfile, err)
	}
	defer f.Close()

	for _, tweet := range tweets {
		log.Printf("%v - %v\n", tweet.IdStr, tweet.User)
		f.WriteString(fmt.Sprintf("%v,%v,%v,%v\n", timestamp, tweet.User.ScreenName, tweet.RetweetCount, tweet.IdStr))

	}

}
