package main

import (
	"errors"
	"fmt"
)

		

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
    handler, ok := c.handlers[cmd.name]
    if !ok {
        return fmt.Errorf("unknown command: %s", cmd.name)
    }

    return handler(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return errors.New("username required")
    }

    err := s.config.SetUser(cmd.args[0])
    if err != nil {
        return err
    }

    fmt.Printf("User has been set to %s\n", cmd.args[0])

    return nil
}
