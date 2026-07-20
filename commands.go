package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Frevens/gator/internal/database"
	"github.com/google/uuid"
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
func handlerAgg(s *state, cmd command) error {
    feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
    if err != nil {
        return err
    }

    fmt.Printf("%+v\n", feed)
    return nil
}
func handlerAddFeed(s *state, cmd command) error {
    if len(cmd.args) < 2 {
        return errors.New("feed name and url required")
    }
    ctx := context.Background()

    user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
    if err != nil {
        return err
    }

    feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
        ID:        uuid.New(),
        Name:      cmd.args[0],
        Url:       cmd.args[1],
        UserID:    user.ID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    })
    if err != nil {
        return err
    }

    fmt.Printf("%+v\n", feed)
    return nil
}

func handlerUsers(s *state, cmd command) error {
    users, err := s.db.GetUsers(context.Background())
    if err != nil {
        return err
    }

    for _, user := range users {
        if user.Name == s.cfg.CurrentUserName {
            fmt.Printf("* %s (current)\n", user.Name)
        } else {
            fmt.Printf("* %s\n", user.Name)
        }
    }

    return nil
}
func handlerReset(s *state, cmd command) error {
    if err := s.db.DeleteAllUsers(context.Background()); err != nil {
        return err
    }

    fmt.Println("All users have been deleted.")
    return nil
}

func handlerLogin(s *state, cmd command) error {
    //consultar con query GetUser
    //si falla porque no existe, devolver error
    //solo si existe, hacer SetUser
   if len(cmd.args) == 0 {
        return errors.New("username required")
    }

    user, err := s.db.GetUser(context.Background(), cmd.args[0])
    if err != nil {
        return errors.New("user not found")
    }

    err = s.cfg.SetUser(user.Name)
    if err != nil {
        return err
    }

    fmt.Printf("User %s has been set as current user\n", user.Name)
    return nil
}
func handlerRegister(s *state, cmd command) error {
    // validar args
    if len(cmd.args) == 0 {
        return errors.New("username required")
    }

   
    // armar params
    _, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
        ID:        uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name:      cmd.args[0],
    })
    if err != nil {
        return err
    }

    // guardar usuario actual en config
    err = s.cfg.SetUser(cmd.args[0])
    if err != nil {
        return err
    }

    // imprimir algo útil
    fmt.Printf("User %s has been registered and set as current user\n", cmd.args[0])

    return nil
}       
   