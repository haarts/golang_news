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

	pollFeeds()
}

func pollFeeds() {
	blogItems := make(chan *rss.Item)
	hnItems := make(chan *rss.Item)
	redditItems := make(chan *rss.Item)

	go poller("https://blog.golang.org/feed.atom", blogItems)
	go poller("https://news.ycombinator.com/rss", hnItems)
	go poller("https://www.reddit.com/r/golang.rss", redditItems)

	for {
		select {
		case item := <-blogItems:
			blogItem(item)
		case item := <-hnItems:
			hnItem(item)
		case item := <-redditItems:
			redditItem(item)
		}
	}
}

func hnItem(item *rss.Item) {
	if match, _ := regexp.MatchString(`\w Go( |$|\.)`, item.Title); match {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		PostTweet(short_title + " " + item.Link + " #hackernews")
	} else {
		log.Printf("Ignoring Hackernews item: %s", item.Title)
	}
}

func blogItem(item *rss.Item) {
	short_title := item.Title
	if len(short_title) > 100 {
		short_title = short_title[:99] + "…"
	}
	PostTweet(short_title + " " + item.Link + " #go_blog")
}

func redditItem(item *rss.Item) {
	re := regexp.MustCompile(`([^"]+)">\[link\]`)
	matches := re.FindStringSubmatch(item.Content)
	if len(matches) == 2 {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		PostTweet(short_title + " " + matches[1] + " #reddit")
	}
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
