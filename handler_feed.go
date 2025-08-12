package main

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/flixflop/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 2 {
        return fmt.Errorf("ADDFEED expects two arguments.")
    }
    name := cmd.args[0]
    url := cmd.args[1]
    feed, err := s.db.CreateFeed(context.Background(),
        database.CreateFeedParams{
            ID:        uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Name:      name, 
            Url:       url,
            UserID:    user.ID,
        })
    if err != nil {
        return fmt.Errorf("Could not create a feed: %w", err)
    }

    _, err = s.db.CreateFeedFollow(context.Background(), 
        database.CreateFeedFollowParams{
        ID:        uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID:    user.ID,
        FeedID:    feed.ID,
    })
    if err != nil {
        return fmt.Errorf("Could not create a feed follow assocation: %w", err)
    }

    rssFeed, err := fetchFeed(context.Background(), url)
    if err != nil {
        return fmt.Errorf("Error fetching RSS Feed from %s: %w", url, err)
    }
    fmt.Println(rssFeed)
    return nil
}

func handlerFeeds(s *state, cmd command) error {
    if len(cmd.args) > 0 {
        return fmt.Errorf("FEEDS does not expect any arguments, however %d were given", len(cmd.args))
    }
    feeds, err := s.db.GetFeeds(context.Background())
    if err != nil {
        return fmt.Errorf("Error getting all feeds: %w", err)
    }
    for _, feed := range feeds {
        printFeed(feed)
    }
    return nil 
}

func printFeed(feed database.GetFeedsRow) {
    fmt.Printf(" * Name:        %v\n", feed.Name)
    fmt.Printf(" * URL:         %v\n", feed.Url)
    fmt.Printf(" * Username:    %v\n", feed.UserName)
}

