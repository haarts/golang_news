package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	rss "github.com/jteeuwen/go-pkg-rss"
	"net/url"
	"os"
	"time"
	"log"
	"regexp"
)

const timeout = 50

func main() {
	log.SetOutput(os.Stdout)

	//go PollFeed("http://blog.golang.org/feed.atom", itemHandlerGoBlog)
	//go PollFeed("https://news.ycombinator.com/rss", itemHandlerHackerNews)
	PollFeed("http://www.reddit.com/r/golang.rss", itemHandlerReddit)
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
		if match, _ := regexp.MatchString(`\w Go`, item.Title); match {
			fmt.Println(item)
		}
	}

	genericItemHandler(feed, ch, newItems, f)
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
	anaconda.SetConsumerKey("9AuRCuPxHAuvPdLIt9Sg")
	anaconda.SetConsumerSecret("14pYoR2B3IAeRsQiOZrSrjweseNj8YW97XbHjwkp4")
	api := anaconda.NewTwitterApi("404687248-c3I6DzhE3KpZvX7lX1J3DFw1AFrtzj8xiwldvBJG", "P7pFC7ZEpKp5WvqtK6xnY0GPyH8Fqt9eOjDrjqi0jlIbO")

	v := url.Values{}
	api.PostTweet(tweet, v)
}
