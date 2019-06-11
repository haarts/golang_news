package feed

import (
	"github.com/SlyMarbo/rss"
)

type Feed struct {
	URL         string
	ItemHandler func(*rss.Item) string
}
