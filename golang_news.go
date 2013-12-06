package main

import (
	"bytes"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	rss "github.com/haarts/go-pkg-rss"
	"log"
	"net/url"
	"os"
	"regexp"
	"time"
)

const timeout = 50

var firstHN = true
var firstReddit = true
var firstGoBlog = true

func main() {
	log.SetOutput(os.Stdout)

	//go PollFeed("http://blog.golang.org/feed.atom", itemHandlerGoBlog)
	PollFeed("https://news.ycombinator.com/rss", itemHandlerHackerNews)
	//PollFeed("http://www.reddit.com/r/golang.rss", itemHandlerReddit)
}

func PollFeed(uri string, itemHandler rss.ItemHandler) {
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
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
				buffer := bytes.NewBufferString(short_title)
				buffer.Truncate(100)
				short_title = buffer.String()
				short_title = short_title + "…"
			}
			PostTweet(short_title + " " + item.Links[0].Href + " #hackernews")
		}
	}

	if firstHN {
		firstHN = false
	} else {
		genericItemHandler(feed, ch, newItems, f)
	}
}

func itemHandlerGoBlog(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		fmt.Println(item)
	}

	genericItemHandler(feed, ch, newItems, f)
}

func itemHandlerReddit(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		re := regexp.MustCompile(`([^"]+)">\[link\]`)
		matches := re.FindStringSubmatch(item.Description)
		if len(matches) == 2 {
			fmt.Println(matches[1])
		}
	}

	//Actually I don't think passing everything around is necessary as the closure remembers context
	genericItemHandler(feed, ch, newItems, f)
}

func PostTweet(tweet string) {
	anaconda.SetConsumerKey(ReadConsumerKey())
	anaconda.SetConsumerSecret(ReadConsumerSecret())
	api := anaconda.NewTwitterApi(ReadAccessToken(), ReadAccessTokenSecret())

	v := url.Values{}
	_, err := api.PostTweet(tweet, v)
	if err != nil {
		log.Printf("Error posting tweet: %s", err)
	}
}
