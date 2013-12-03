package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	rss "github.com/jteeuwen/go-pkg-rss"
	"net/url"
	"os"
	"time"
)

const timeout = 50

func main() {
	go PollFeed("http://blog.golang.org/feed.atom", timeout)
	go PollFeed("https://news.ycombinator.com/rss", timeout)
	PollFeed("http://www.reddit.com/r/golang.rss", timeout)
}

func PollFeed(uri string, timeout int) {
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
	fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {
	fmt.Printf("%d new item(s) in %s\n", len(newItems), feed.Url)
}

func PostTweet(tweet string) {
	anaconda.SetConsumerKey("9AuRCuPxHAuvPdLIt9Sg")
	anaconda.SetConsumerSecret("14pYoR2B3IAeRsQiOZrSrjweseNj8YW97XbHjwkp4")
	api := anaconda.NewTwitterApi("404687248-c3I6DzhE3KpZvX7lX1J3DFw1AFrtzj8xiwldvBJG", "P7pFC7ZEpKp5WvqtK6xnY0GPyH8Fqt9eOjDrjqi0jlIbO")

	v := url.Values{}
	api.PostTweet(tweet, v)
}
