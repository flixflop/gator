package main

import (
    "context"
    "errors"
    "log"
    "fmt"
    "time"
    "strings"
    "database/sql"
 
    "github.com/google/uuid"
    "github.com/flixflop/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
    log.Println("handlerAgg: received command arguments", cmd.args)
    if len(cmd.args) != 1 {
        return errors.New("AGG takes one argument.")
    }
    durationString := cmd.args[0]
    timeBetweenRequests, err := time.ParseDuration(durationString)
    if err != nil {
        return fmt.Errorf("Could not parse time_between_requests argument %s: %w", durationString, err)
    }
    fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

    ticker := time.NewTicker(timeBetweenRequests)
    for ; ; <-ticker.C {
            scrapeFeeds(s)
    }
   return nil
}

func scrapeFeeds(s *state) {
    feed, err := s.db.GetNextFeedToFetch(context.Background())
    if err != nil {
        log.Println("Could not get next feed to fetch!")
    }
    log.Println("Found a feed to fetch!")
    scrapeFeed(s, feed)
}

func scrapeFeed(s *state, feed database.Feed) {
    _, err := s.db.MarkFeedFetched(context.Background(), feed.ID)
    if err != nil {
        log.Printf("Could not mark feed %s as 'fetched': %w", feed.Name, err)
    }

    feedData, err := fetchFeed(context.Background(), feed.Url)
    if err != nil {
        log.Printf("Could not collect feed %s: %w", feed.Name, err)
    }
    for _, item := range feedData.Channel.Item {
        publishedAt := sql.NullTime{}
        if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
            publishedAt = sql.NullTime{
                Time: t,
                Valid: true,
            }
        }
        _, err = s.db.CreatePost(context.Background(),
            database.CreatePostParams{
                ID:          uuid.New(),
                CreatedAt:   time.Now(),
                UpdatedAt:   time.Now(),
                Title:       item.Title,
                Url:         item.Link,
                Description: sql.NullString{
                    String: item.Description,
                    Valid:  true,
                },
                PublishedAt: publishedAt,
                FeedID:      feed.ID,
        })
        if err != nil {
            if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
                continue
            }
            log.Printf("Couldn't create post: %v", err)
            continue
        }
    }
    log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

