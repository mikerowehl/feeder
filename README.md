# feeder

[![Go](https://github.com/mikerowehl/feeder/actions/workflows/go.yml/badge.svg)](https://github.com/mikerowehl/feeder/actions/workflows/go.yml)

# Feeder - command line syndicated feed tool

<img src="gopher_bowl.png" width="200" alt="Go Gopher with a bowl of RSS feed icons">

`feeder` is a simple command line tool for tracking, fetching, and reading syndicated content.
It's built on top of the fantastic [gofeed package](https://github.com/mmcdole/gofeed) so it supports
a bunch of syndication formats. The list of feeds to keep track of and their items are stored in a
minimal sqlite database locally, so there's no other services to setup or any database or web service
to configure. When you want to read the batch of new content it gets written into a local file, all 
self contained and on a single page.

## Example Usage

Add a few feeds to the database:

* `feeder add https://rowehl.com/feed.xml`
* `feeder add https://pluralistic.net/feed/`

Pull down the content:

* `feeder fetch`

Make a page with links to each of the items from the feeds:

* `feeder read`

Right now the read command just outputs to a file named feeder.html in the current directory.
