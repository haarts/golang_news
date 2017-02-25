package main

import (
	"log"
	"net/http"
	"time"

	"github.com/SlyMarbo/rss"
)

func Poller(uri string, items chan<- *rss.Item) {
	feed := getFeed(uri)

	for _, item := range feed.Items {
		// ignore the first batch of items
		//items <- item
	}

	<-time.After(time.Duration(feed.Refresh.Sub(time.Now())))

	for {
		err := feed.Update()
		if err != nil {
			log.Printf("Error fetching %s: %s", uri, err)
		} else {
			for _, item := range feed.Items {
				items <- item
			}
		}

		<-time.After(time.Duration(feed.Refresh.Sub(time.Now())))
	}
}

func getFeed(uri string) *rss.Feed {
	f := func(url string) (resp *http.Response, err error) {
		client := http.DefaultClient
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", "@golang_news twitter bot")

		return client.Do(req)
	}

	for {
		feed, err := rss.FetchByFunc(f, uri)
		if err == nil {
			return feed
		}
		log.Printf("Error fetching %s: %s", uri, err)

		<-time.After(time.Duration(50 * time.Second))
	}
}
