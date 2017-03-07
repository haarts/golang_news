package main

import "C"
import (
	"log"
	"regexp"

	"github.com/SlyMarbo/rss"
)

type Feed struct {
	URL         string
	ItemHandler func(*rss.Item) string
}

func List() []Feed {
	return []Feed{
		Feed{
			URL:         "http://www.juliabloggers.com/feed/",
			ItemHandler: blogItem,
		},
		Feed{
			URL:         "https://news.ycombinator.com/rss",
			ItemHandler: hnItem,
		},
		Feed{
			URL:         "https://www.reddit.com/r/julia.rss",
			ItemHandler: redditItem,
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
