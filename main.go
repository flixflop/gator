package main

import (
    "log"
    "os"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/flixflop/gator/internal/config"
    "github.com/flixflop/gator/internal/database"
)

type state struct {
    config *config.Config
    db     *database.Queries
}

func main() {
    cfg, err := config.Read()
    if err != nil {
        log.Fatal(err)
    }

    // connect to DB
    db, err := sql.Open("postgres", cfg.DBURL)
    dbQueries := database.New(db)
    // set state from read config
    s := state{
        config: &cfg,
        db:     dbQueries,
    }
    // initialize and register command line arguments
    c := commands{
        mapCmdToFunc: make(map[string]func(*state, command) error),
    }
    c.register("login", handlerLogin)
    c.register("register", handlerRegister)
    c.register("reset", handlerReset)
    c.register("users", handlerUsers)
    c.register("agg", handlerAgg)
    c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
    c.register("feeds", handlerFeeds)
    c.register("follow", middlewareLoggedIn(handlerFollow))
    c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
    c.register("following", middlewareLoggedIn(handlerFollowing))
    c.register("browse", middlewareLoggedIn(handlerBrowse))
 
    // get actual command line arguments
    args := os.Args
    if len(args) < 2 {
        log.Fatal("Expected at least a single argument")
    }
    cmd := command{
        name: args[1],
        args: args[2:],
    }
    err = c.run(&s, cmd)
    log.Println("Received an error:", err)
    if err != nil {
        log.Fatal(err)
    }
}
