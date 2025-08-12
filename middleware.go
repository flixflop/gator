package main

import (
    "fmt"
    "context"
    "github.com/flixflop/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
   return func(s *state, cmd command) error {
        user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
        if err != nil {
            return fmt.Errorf("Error getting user '%s' from database: %w", s.config.CurrentUserName, err)
        }
        return handler(s, cmd, user)
    }
}
