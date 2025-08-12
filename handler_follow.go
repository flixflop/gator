package main

import (
    "fmt"
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/flixflop/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("FOLLOW takes a single argument") 
    }
    url := cmd.args[0]

    feed, err := s.db.GetFeedByUrl(context.Background(), url)
    if err != nil {
        return fmt.Errorf("Could not find a feed with url %s: %w", url, err)
    }

    feedFollow, err := s.db.CreateFeedFollow(context.Background(), 
        database.CreateFeedFollowParams{
            ID: uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            UserID:    user.ID,
            FeedID:    feed.ID,
        })
    if err != nil {
        return fmt.Errorf("Could not create a feed follow assocation: %w", err)
    }
    fmt.Printf(" * Name:   %v\n", feedFollow.UserName)
    fmt.Printf(" * Feed:   %v\n", feedFollow.FeedName)
    return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
    feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {
        return fmt.Errorf("Error getting all followed feeds for user %s: %v", user.Name, err)
    }
    for _, feedFollow := range feedFollows {
        fmt.Printf("%s", feedFollow.FeedName) 
    }
    return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("UNFOLLOW expects one argument.")
    }

    url := cmd.args[0]
    err := s.db.DeleteFeedFollow(context.Background(),
        database.DeleteFeedFollowParams{
            Name: user.Name,
            Url:  url,
        })
    if err != nil {
        return fmt.Errorf("Could not delete followed feed: %w", err)
    }
    fmt.Println("Successfully unfollowed!")
    return nil
}

