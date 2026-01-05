<img src="gopher_bowl.png" width="200" alt="Go Gopher with a bowl of RSS feed icons">

# Feeder - command line feed tool

[![Go](https://github.com/mikerowehl/feeder/actions/workflows/go.yml/badge.svg)](https://github.com/mikerowehl/feeder/actions/workflows/go.yml)

`feeder` is a simple command line tool for tracking, fetching, and reading syndicated content.
It's built on top of the fantastic [gofeed package](https://github.com/mmcdole/gofeed) so it supports
a bunch of syndication formats. The list of feeds to keep track of and their items are stored in a
minimal sqlite database locally, so there's no other services to setup or any database or web service
to configure. When you want to read the batch of new content it gets written into a local file, all 
self contained and on a single page.

## Example Usage

Add a few feeds to the database:

```
feeder add https://www.youtube.com/@tested
feeder add https://rowehl.com/feed.xml
feeder add https://pluralistic.net/feed/
feeder add https://hackaday.com/blog/feed/
```

The daily command fetches the most recent feed content, writes any unread items into an HTML file
in the current directory, marks all the items as read in the database, and attempts to run a
command that will open the generated file in the default browser.

```feeder daily```

