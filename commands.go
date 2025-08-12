package main

import "errors"

type command struct {
    name string
    args []string
}

type commands struct {
    mapCmdToFunc map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
    f, ok := c.mapCmdToFunc[cmd.name]
    if !ok {
        return errors.New("command " + cmd.name + " could not be found.")
    }

    err := f(s, cmd)
    if err != nil {
        return err
    }
    return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.mapCmdToFunc[name] = f
}

