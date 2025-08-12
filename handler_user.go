package main

import (
    "fmt"
    "context"
    "log"
    "strings"
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/flixflop/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
    log.Println("handlerLogin: Got supplied", len(cmd.args), "command line arguments: ", cmd.args)
    if len(cmd.args) == 0 {
        log.Println("handlerLogin: I am reaching the error return")
        return errors.New("Expecting command arguments, none were supplied!")
    }
    username := strings.Join(cmd.args, " ")

    // check that user is already in DB
    _, err := s.db.GetUser(context.Background(), username)
    if err != nil {
        return errors.New("User " + username + " does not seem to exist.")
    }
    s.config.SetUser(username)
    fmt.Println("User", username, "has been set.")
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return errors.New("Expecting a single argument.")
    }
    name := cmd.args[0]
    user, err := s.db.CreateUser(context.Background(), 
        database.CreateUserParams{
            ID:        uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Name:      name, 
        })
    if err != nil {
        return errors.New("User " + name + " already exists")
    }
    logUser(user)
    s.config.SetUser(name)
   return nil
}

func handlerUsers(s *state, cmd command) error {
    if len(cmd.args) > 1 {
        return errors.New("USERS does not take arguments.")
    }
    users, err := s.db.GetUsers(context.Background())
    if err != nil {
        return errors.New("Could not get all users.")
    }

    for _, user := range users {
        if user.Name != s.config.CurrentUserName {
            fmt.Println("* ", user.Name)
        } else {
            fmt.Println("* ", user.Name, "(current)")
        }
    }
    return nil
}

func logUser(user database.User) {
    log.Println("User", user.Name, "has been created.")
    log.Println("UUID", user.ID)
    log.Println("CreatedAt", user.CreatedAt)
    log.Println("UpdatedAt", user.UpdatedAt)
}
