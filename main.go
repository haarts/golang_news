package main

import (
	"log"
	"os"
	"plugin"
	"reflect"
	"github.com/haarts/golang_news/feed"
	"github.com/SlyMarbo/rss"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting")

	feeds, err := feeds()
	if err != nil {
		return
	}
	tweets := make(chan string)

	go pollFeeds(tweets, feeds)
	postTweets(tweets)
}

func feeds() ([]feed.Feed, error) {
	p, err := plugin.Open(os.Args[1])
	if err != nil {
		log.Printf("Error reading plugin: %s", err)
		return nil, err
	}
	listFunc, err := p.Lookup("List")
	if err != nil {
		log.Printf("Error looking up 'List': %s", err)
		return nil, err
	}

	return listFunc.(func() []feed.Feed)(), nil
}

func pollFeeds(publishTweet chan string, feeds []feed.Feed) {
	itemProducers := []chan *rss.Item{}
	for _, feed := range feeds {
		itemProducer := make(chan *rss.Item)
		itemProducers = append(itemProducers, itemProducer)
		log.Printf("Adding '%s'", feed.URL)
		go poller(feed.URL, itemProducer)
	}

	cases := make([]reflect.SelectCase, len(itemProducers))
	for i, ch := range itemProducers {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}

	activeItemProducers := len(cases)
	for activeItemProducers > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			activeItemProducers -= 1
			continue
		}

		publishTweet <- feeds[chosen].ItemHandler(value.Interface().(*rss.Item))
	}
}

func postTweets(tweets chan string) {
	for tweet := range tweets {
		if tweet != "" {
			PostTweet(tweet)
		}
	}
}
