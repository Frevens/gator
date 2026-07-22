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

func handlerFollow(s *state, cmd command, user database.User) error {
    // 1. Define context first so it's available for all DB calls
    ctx := context.Background()

    // 2. Validate input
    if len(cmd.args) != 1 {
        return errors.New("follow requires a feed URL")
    }
    url := cmd.args[0]

    // 3. Get the feed
    feed, err := s.db.GetFeedByURL(ctx, url)
    if err != nil {
        return err
    }

    // 4. Current user managed by a middleware


    // 5. Create the follow record
    // Note: CreateFeedFollowParams is a struct TYPE. We create an INSTANCE using {...}
    follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
    ID:        uuid.New(),
    UserID:    user.ID,
    FeedID:    feed.ID,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
})
    if err != nil {
        return err
    }

    // 6. Success output
    // Use %+v to print struct fields, or access specific fields like feed.Name
    fmt.Printf("Followed Feed: %s\nUser: %s\n", follow.FeedName, follow.UserName)
    
    return nil
}   

func handlerFollowing(s *state, cmd command, user database.User) error {
    ctx := context.Background()

    // Current user managed by a middleware

    follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
    if err != nil {
        return err
    }

    for _, follow := range follows {
        fmt.Println(follow.FeedName)
    }

    return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
    // 1. Define context first so it's available for all DB calls
    ctx := context.Background()

    // 2. Validate input
    if len(cmd.args) != 1 {
        return errors.New("follow requires a feed URL")
    }
    url := cmd.args[0]

    // 3. Get the feed
    feed, err := s.db.GetFeedByURL(ctx, url)
    if err != nil {
        return err
    }
    // 4. Current user managed by a middleware
    // 5. Delete the follow record
    err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
        UserID: user.ID,
        FeedID: feed.ID,
    })
    if err != nil {
        return err
    }

    fmt.Printf("User: %s has unfollowed Feed: %s\n", user.Name, feed.Name)
    
    return nil

}

func handlerAgg(s *state, cmd command) error {
    feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
    if err != nil {
        return err
    }

    fmt.Printf("%+v\n", feed)
    return nil
}
func handlerAddFeed(s *state, cmd command, user database.User) error {
    if len(cmd.args) < 2 {
        return errors.New("feed name and url required")
    }
    ctx := context.Background()

    // Current user managed by a middleware

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

    _, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
        ID:        uuid.New(),
        UserID:    user.ID,
        FeedID:    feed.ID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    })
    if err != nil {
        return err
    }

    fmt.Printf("%+v\n", feed)
    return nil
}

func handlerFeeds(s *state, cmd command) error {
    feeds, err := s.db.ListFeeds(context.Background())
    if err != nil {
        return err
    }
    for _,feed := range feeds {
        fmt.Printf("Feed: %s\nURL: %s\nCreated by: %s\n\n",
        feed.FeedName,
        feed.FeedUrl,
        feed.UserName,
        )
    }
    
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
   