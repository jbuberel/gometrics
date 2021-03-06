package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	cls "github.com/jbuberel/gometrics/gerritcls"
	g "github.com/jbuberel/gometrics/github"
	iss "github.com/jbuberel/gometrics/gitissues"
	mobile "github.com/jbuberel/gometrics/mobile"
	r "github.com/jbuberel/gometrics/reddit"
	t "github.com/jbuberel/gometrics/twitter"
	tw "github.com/jbuberel/gometrics/twittermobile"
)

var dirname = flag.String("dir", "/usr/local/google/home/jbuberel/gometrics", "The directory where the file will be written")
var since = flag.String("since", time.Now().Add(-47*time.Hour).Format("2006-01-02"), "The start date of the search window, in YYYY-MM-DD format.")
var until = flag.String("until", time.Now().Add(-23*time.Hour).Format("2006-01-02"), "The end date of the search window, in YYYY-MM-DD format.")

var clsToggle = flag.Bool("cls", true, "Set to false to disable Gerrit CLs capture")
var githubToggle = flag.Bool("github", true, "set to false to disable capture of Github stats")
var issuesToggle = flag.Bool("issues", true, "set to false to disable capture of Go issues")
var redditToggle = flag.Bool("reddit", true, "set to false to disable capture of reddit data")
var twitterToggle = flag.Bool("twitter", true, "set to false to disable capture of twitter data")
var mobileToggle = flag.Bool("mobile", true, "set to false to disable capture of mobile import data")
var twittermobileToggle = flag.Bool("twittermobile", true, "set to false to disable capture of #golang #mobile data")

var twitterConsumerKey string
var twitterConsumerSecret string
var twitterAccessToken string
var twitterSecretToken string
var githubClientId string
var githubSecretKey string
var githubSecretToken string

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	log.Printf("Looking through env vars\n")
	for _, e := range os.Environ() {
		parts := strings.Split(e, "=")
		if len(parts) == 2 {
			if parts[0] == "twitter_consumer_key" {
				twitterConsumerKey = string(parts[1])
				log.Printf("twitter_consumer_key set from environ to: %v\n", twitterConsumerKey)
			} else if parts[0] == "twitter_consumer_secret" {
				twitterConsumerSecret = string(parts[1])
				log.Printf("twitter_consumer_secret set from environ to: %v\n", twitterConsumerSecret)
			} else if parts[0] == "twitter_access_token" {
				twitterAccessToken = string(parts[1])
				log.Printf("twitter_access_token set from environ to: %v\n", twitterAccessToken)
			} else if parts[0] == "twitter_secret_token" {
				twitterSecretToken = string(parts[1])
				log.Printf("twitter_secret_token set from environ to: %v\n", twitterSecretToken)
			} else if parts[0] == "github_client_id" {
				githubClientId = string(parts[1])
				log.Printf("github_client_id set from environ to: %v\n", githubClientId)
			} else if parts[0] == "github_secret_key" {
				githubSecretKey = string(parts[1])
				log.Printf("github_secret_key set from environ to: %v\n", githubSecretKey)
			} else if parts[0] == "github_secret_token" {
				githubSecretToken = string(parts[1])
				log.Printf("github_secret_key set from environ to: %v\n", githubSecretToken)
			}
		}

	}

	if len(twitterConsumerSecret) == 0 || len(twitterConsumerKey) == 0 || len(twitterSecretToken) == 0 {
		log.Println("Unable to obtain twitter keys from environment variables!!")
		os.Exit(1)
	}
	if len(githubClientId) == 0 || len(githubSecretKey) == 0 || len(githubSecretToken) == 0 {
		log.Println("Unable to obtain github keys from environment variables!!")
		os.Exit(1)
	}

}

func main() {
	log.Printf("Saving data to directory: %v\n", *dirname)
	log.Println("Starting gometrics capture")
	if *redditToggle {
		r.Capture(dirname)
	}
	if *twitterToggle {
		t.Capture(dirname, since, until, twitterConsumerKey, twitterConsumerSecret, twitterAccessToken, twitterSecretToken)
	}
	if *githubToggle {
		g.Capture(dirname, githubClientId, githubSecretKey, githubSecretToken)
	}
	if *clsToggle {
		cls.Capture(dirname)
	}
	if *issuesToggle {
		iss.Capture(dirname, githubClientId, githubSecretKey, githubSecretToken)
	}
	if *mobileToggle {
		mobile.Capture(dirname)
	}
	if *twittermobileToggle {
		tw.Capture(dirname, since, until, twitterConsumerKey, twitterConsumerSecret, twitterAccessToken, twitterSecretToken)
	}
	log.Println("gometrics capture complete")
}
