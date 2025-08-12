package main

import (
    "context"
    "fmt"
)

func handlerReset(s *state, cmd command) error {
    if len(cmd.args) > 1 {
        return fmt.Errorf("RESET does not take arguments.")
    }
    err := s.db.DeleteUsers(context.Background())
    if err != nil {
        return fmt.Errorf("Could not delete all users: %w", err)
    }
    fmt.Println("Successfully deleted all users")
    return nil
}

