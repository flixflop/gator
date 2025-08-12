package main

import (
    "fmt"
    "context"
    "strconv"

    "github.com/flixflop/gator/internal/database"
)
 
func handlerBrowse(s *state, cmd command, user database.User) error {
    if len(cmd.args) > 1 {
        return fmt.Errorf("BROWSE takes at most one argument; %d were given", len(cmd.args))
    }

    limit := 2
    if len(cmd.args) == 1 {
        specifiedLimit, err := strconv.Atoi(cmd.args[0]); if err == nil {
            limit = specifiedLimit
        } else {
            return fmt.Errorf("BROWSE got an invalid limit: %w", err)
        }
    }
    posts, err := s.db.GetPostByUser(context.Background(), 
        database.GetPostByUserParams{
            Name:  user.Name,
            Limit: int32(limit),
    })
    if err != nil {
        return fmt.Errorf("Error getting posts for user %s: %w", user.Name, err)
    }
    for _, post := range posts {
        fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
        fmt.Printf("--- %s ---\n", post.Title)
        fmt.Printf("    %v\n", post.Description.String)
        fmt.Printf("Link: %s\n", post.Url)
        fmt.Println("=====================================")
    }
    return nil
}
 
