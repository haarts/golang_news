package main

import (
	"log"
	"net/url"
	"os"
	"regexp"

	"github.com/ChimeraCoder/anaconda"
	"github.com/SlyMarbo/rss"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting")
	tweets := make(chan string)

	go pollFeeds(tweets)
	postTweets(tweets)
}

func pollFeeds(publishTweet chan string) {
	// variable number of select/case: https://play.golang.org/p/8zwvSk4kjx
	blogItems := make(chan *rss.Item)
	hnItems := make(chan *rss.Item)
	redditItems := make(chan *rss.Item)

	go poller("https://blog.golang.org/feed.atom", blogItems)
	go poller("https://news.ycombinator.com/rss", hnItems)
	go poller("https://www.reddit.com/r/golang.rss", redditItems)

	for {
		select {
		case item := <-blogItems:
			publishTweet <- blogItem(item)
		case item := <-hnItems:
			publishTweet <- hnItem(item)
		case item := <-redditItems:
			publishTweet <- redditItem(item)
		}
	}
}

func postTweets(tweets chan string) {
	for tweet := range tweets {
		if tweet != "" {
			PostTweet(tweet)
		}
	}
}

func hnItem(item *rss.Item) string {
	if match, _ := regexp.MatchString(`\w Go( |$|\.)`, item.Title); match {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		return short_title + " " + item.Link + " #hackernews"
	} else {
		log.Printf("Ignoring Hackernews item: %s", item.Title)
		return ""
	}
}

func blogItem(item *rss.Item) string {
	short_title := item.Title
	if len(short_title) > 100 {
		short_title = short_title[:99] + "…"
	}
	return short_title + " " + "https:" + item.Link + " #go_blog"
}

func redditItem(item *rss.Item) string {
	re := regexp.MustCompile(`([^"]+)">\[link\]`)
	matches := re.FindStringSubmatch(item.Content)
	if len(matches) == 2 {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		return short_title + " " + matches[1] + " #reddit"
	}
	return ""
}

func PostTweet(tweet string) {
	anaconda.SetConsumerKey(ReadConsumerKey())
	anaconda.SetConsumerSecret(ReadConsumerSecret())
	api := anaconda.NewTwitterApi(ReadAccessToken(), ReadAccessTokenSecret())

	log.Printf("Posting tweet: %s", tweet)

	v := url.Values{}
	_, err := api.PostTweet(tweet, v)
	if err != nil {
		log.Printf("Error posting tweet: %s", err)
	}
}
