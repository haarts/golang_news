package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/ChimeraCoder/anaconda"
	rss "github.com/jteeuwen/go-pkg-rss"
)

const timeout = 50

var first = map[string]bool{}

//var first = map[string]bool{"reddit": true, "hn": true, "go": true}

type hn struct{}
type blog struct{}
type reddit struct{}

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting")

	go PollFeed("http://blog.golang.org/feed.atom", blog{})
	go PollFeed("https://news.ycombinator.com/rss", hn{})
	PollFeed("http://www.reddit.com/r/golang.rss", reddit{})
}

func PollFeed(uri string, itemHandler rss.ItemHandler) {
	feed := rss.New(timeout, true, chanHandler, itemHandler.ProcessItems)

	for {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		err := feed.FetchClient(uri, client, nil)
		if err != nil {
			// don't die, just log and retry.
			log.Printf("Error fetching %s: %s", uri, err)
		} else {
			log.Printf("Fetched %s\n", uri)
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

func (h hn) ProcessItems(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
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

func (b blog) ProcessItems(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
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

func (r reddit) ProcessItems(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	f := func(item *rss.Item) {
		re := regexp.MustCompile(`([^"]+)">\[link\]`)
		matches := re.FindStringSubmatch(item.Content.Text)
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
