package main

import "C"
import (
	"log"
	"regexp"

	"github.com/haarts/golang_news/feed"
	"github.com/SlyMarbo/rss"
)



func List() []feed.Feed {
	return []feed.Feed{
		feed.Feed{
			URL:         "https://www.juliabloggers.com/feed/",
			ItemHandler: blogItem,
		},
		feed.Feed{
			URL:         "https://news.ycombinator.com/rss",
			ItemHandler: hnItem,
		},
	}
}

func hnItem(item *rss.Item) string {
	if match, _ := regexp.MatchString(`\w Julia( |$|\.)`, item.Title); match {
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
	return short_title + " " + item.Link + " #juliabloggers"
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
