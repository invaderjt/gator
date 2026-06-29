package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/invaderjt/blog-aggregator/internal/config"
	"github.com/invaderjt/blog-aggregator/internal/database"
	"github.com/invaderjt/blog-aggregator/internal/rss"
)

func initialize() *state {
	s := &state{}
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Could not read config file")
	}
	s.Cfg = &cfg
	return s
}

func updateState(s *state, db *sql.DB) {
	dbQueries := database.New(db)
	s.Db = dbQueries

	var cmds commands
	cmds.Commands = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	input := os.Args
	if len(input) < 2 {
		log.Fatalln("Not enough arguments")
	}
	var cmd command
	cmd.Name = input[1]
	cmd.Args = input[2:]
	err := cmds.run(s, cmd)
	if err != nil {
		fmt.Println(err)
	}

}

type state struct {
	Db  *database.Queries
	Cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	Commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	return c.Commands[cmd.Name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Commands[name] = f
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			log.Fatalln("Invalid current user")
		}
		return handler(s, cmd, currentUser)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		log.Fatalln("Login requires username argument")
	}

	name := cmd.Args[0]
	user, exists := s.Db.GetUser(context.Background(), name)
	if exists != nil {
		log.Fatalf("User with name %s does not exist\n", name)
	}

	err := s.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		log.Fatalln("Register requires username argument")
	}

	ctx := context.Background()
	uuid := uuid.New()
	created_at := time.Now()
	updated_at := created_at
	name := cmd.Args[0]

	if _, exists := s.Db.GetUser(ctx, name); exists == nil {
		log.Fatalf("User with name %s already exists\n", name)
	}

	params := database.CreateUserParams{
		ID:        uuid,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      name,
	}

	_, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		log.Fatalln("Error creating user")
	}

	handlerLogin(s, cmd)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.Db.ResetDB(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	for _, user := range users {
		if user == s.Cfg.CurrentUserName {
			fmt.Println(user + " (current)")
		} else {
			fmt.Println(user)
		}
	}
	return nil

}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	data, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(data)

	return nil
}

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
	if len(cmd.Args) < 2 {
		log.Fatalln("Addfeed requires name and url")
	}

	ctx := context.Background()
	uuid := uuid.New()
	created_at := time.Now()
	updated_at := created_at
	name := cmd.Args[0]
	url := cmd.Args[1]
	user_id := currentUser.ID

	params := database.AddFeedParams{
		ID:        uuid,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      name,
		Url:       url,
		UserID:    user_id,
	}

	feed, err := s.Db.AddFeed(ctx, params)
	if err != nil {
		log.Fatalf("Error adding feed: %v", err)
	}
	cmd.Args[0] = cmd.Args[1]
	err = handlerFollow(s, cmd, currentUser)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	for _, feed := range feeds {
		name, err := s.Db.GetNameFromUUID(context.Background(), feed.UserID)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%v | %v | %v\n", feed.Name, feed.Url, name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.Args) < 1 {
		log.Fatalln("Follow requires url argument")
	}

	desiredFeed, err := s.Db.GetFeedFromURL(context.Background(), cmd.Args[0])
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	uuid := uuid.New()
	created_at := time.Now()
	updated_at := created_at
	user_id := currentUser.ID
	feed_id := desiredFeed.ID

	params := database.CreateFeedFollowParams{
		ID:        uuid,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		UserID:    user_id,
		FeedID:    feed_id,
	}

	_, err = s.Db.CreateFeedFollow(ctx, params)
	if err != nil {
		log.Fatalf("Could not follow %s\n", desiredFeed.Name)
	}

	fmt.Printf("%v is now following %v\n", currentUser.Name, desiredFeed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, currentUser database.User) error {
	following, err := s.Db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		log.Fatalf("Could not retrieve follow list for %s\n", currentUser.Name)
	}

	fmt.Printf("%s is following these feeds:\n", currentUser.Name)
	for _, feed := range following {
		fmt.Printf("%s | %s\n", feed.FeedName, feed.UserName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	toUnfollow, err := s.Db.GetFeedFromURL(context.Background(), cmd.Args[0])
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	user_id := currentUser.ID
	feed_id := toUnfollow.ID

	params := database.UnfollowParams{
		UserID: user_id,
		FeedID: feed_id,
	}
	err = s.Db.Unfollow(ctx, params)
	if err != nil {
		log.Fatalf("Could not unfollow %s", toUnfollow.Name)
	}

	return nil
}
