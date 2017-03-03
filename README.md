# @golang_news
This tiny piece of Go code powers the Golang News account at [Twitter](http://twitter.com/golang_news).

It checks [HackerNews](http://news.ycombinator.com), [Reddit](http://www.reddit.com/r/golang) and the official 
[Go blog](http://blog.golang.org/) for new content.

## Build

`$ go build`

### Plugins

```shell
$ cd plugins/golang_news
$ go build -buildmode=plugin
```

## Run

Add your credential of your personal twitter bot to `settings.json`, look at
`settings.json.example` for, well, an example.

The settings file is expect to be in the same directory as the executable.

Finally run the binary with the first argument pointing to the plugin (`*.so`
file):

```shell
$ ./golang_news plugins/golang_news/golang_news.so
```
