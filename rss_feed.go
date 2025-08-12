package main

import (
    "context"
    "errors"
    "io"
    "encoding/xml"
    "net/http"
    "html"
)

type RSSFeed struct {
    Channel struct {
        Title       string    `xml:"title"`
        Link        string    `xml:"link"`
        Description string    `xml:"description"`
        Item        []RSSItem `xml:"item"`
    } `xml:"channel"`
}

type RSSItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
    client := &http.Client{}
    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
        return &RSSFeed{}, errors.New("Error creating request.")
    }
    req.Header.Set("User-Agent", "gator")
    resp, err := client.Do(req)
    if err != nil {
        return &RSSFeed{}, errors.New("Error performing request.")
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return &RSSFeed{}, errors.New("Error fetching resource from " + feedURL)
    }

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return &RSSFeed{}, errors.New("Error reading the response body.")
    }
    var feed RSSFeed
    err = xml.Unmarshal(data, &feed)
    if err != nil {
        return &RSSFeed{}, errors.New("Error unmarshaling XML.")
    }
    err = htmlUnescapeFeed(&feed)
    if err != nil {
        return &RSSFeed{}, errors.New("Error escaping HTML.")
    }
    return &feed, nil
}

func htmlUnescapeFeed(feed *RSSFeed) error {
    feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
    feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
    for i, rssItem := range feed.Channel.Item {
        feed.Channel.Item[i].Title = html.UnescapeString(rssItem.Title)
        feed.Channel.Item[i].Description = html.UnescapeString(rssItem.Description)
    }
    return nil
}

