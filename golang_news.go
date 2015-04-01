package main

import (
	"crypto/tls"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	rss "github.com/haarts/go-pkg-rss"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

const timeout = 50

var first = map[string]bool{}

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting")

	go PollFeed("http://blog.golang.org/feed.atom", itemHandlerGoBlog)
	go PollFeed("https://news.ycombinator.com/rss", itemHandlerHackerNews)
	PollFeed("http://www.reddit.com/r/golang.rss", itemHandlerReddit)
}

func PollFeed(uri string, itemHandler rss.ItemHandler) {
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		retries := 0
	Retry:
		log.Printf("Trying (%d) %s", retries, uri)
		if err := feed.FetchClient(uri, client, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", time.Now().Format(time.RFC3339), uri, err)
			retries += 1
			if retries < 3 {
				goto Retry
			} else {
				return
			}
		}

		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	//noop
}

func genericItemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item, individualItemHandler func(*rss.Item)) {
	log.Printf("%d new item(s) in %s\n", len(newItems), feed.Url)
	for _, item := range newItems {
		individualItemHandler(item)
	}
}

func itemHandlerHackerNews(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		if match, _ := regexp.MatchString(`\w Go( |$|\.)`, item.Title); match {
			short_title := item.Title
			if len(short_title) > 100 {
				short_title = short_title[:99] + "…"
			}
			PostTweet(short_title + " " + item.Links[0].Href + " #hackernews")
		} else {
			log.Printf("Ignoring Hackernews item: %s", item.Title)
		}
	}

	if _, ok := first["hn"]; !ok {
		log.Println("Ignoring first batch of Hackernews")
		first["hn"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func itemHandlerGoBlog(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		short_title := item.Title
		if len(short_title) > 100 {
			short_title = short_title[:99] + "…"
		}
		PostTweet(short_title + " " + item.Links[0].Href + " #go_blog")
	}

	if _, ok := first["go"]; !ok {
		log.Println("Ignoring first batch of Go blog")
		first["go"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func itemHandlerReddit(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		re := regexp.MustCompile(`([^"]+)">\[link\]`)
		matches := re.FindStringSubmatch(item.Description)
		if len(matches) == 2 {
			short_title := item.Title
			if len(short_title) > 100 {
				short_title = short_title[:99] + "…"
			}
			PostTweet(short_title + " " + matches[1] + " #reddit")
		}
	}

	if _, ok := first["reddit"]; !ok {
		log.Println("Ignoring first batch of Reddit")
		first["reddit"] = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func PostTweet(tweet string) {
	anaconda.SetConsumerKey(ReadConsumerKey())
	anaconda.SetConsumerSecret(ReadConsumerSecret())
	api := anaconda.NewTwitterApi(ReadAccessToken(), ReadAccessTokenSecret())

	log.Printf("Posting tweet: %s", tweet[len(tweet)-25:])

	v := url.Values{}
	_, err := api.PostTweet(tweet, v)
	if err != nil {
		log.Printf("Error posting tweet: %s", err)
	}
}
