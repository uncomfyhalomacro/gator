package cli

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/uncomfyhalomacro/gator/internal/config"
	"github.com/uncomfyhalomacro/gator/internal/database"
	"github.com/uncomfyhalomacro/gator/internal/rss"
	"log"
	"time"
)

type State struct {
	Db       *database.Queries
	Config_p *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	FuncFromCommand map[string]func(*State, Command) error
}

func Initialise(s *State) Commands {
	mapping := make(map[string]func(*State, Command) error)
	c := Commands{}
	c.FuncFromCommand = mapping
	c.registerCommand("login", handlerLogin)
	c.registerCommand("register", handlerRegister)
	c.registerCommand("reset", handlerReset)
	c.registerCommand("users", handlerGetUsers)
	c.registerCommand("agg", handlerAggregator)
	c.registerCommand("addfeed", middlewareLoggedIn(s, handlerAddFeed))
	c.registerCommand("follow", middlewareLoggedIn(s, handlerFollow))
	c.registerCommand("following", middlewareLoggedIn(s, handlerFollowing))
	c.registerCommand("feeds", handlerGetFeeds)
	return c
}

func (c *Commands) Run(s *State, cmd Command) error {
	r, ok := c.FuncFromCommand[cmd.Name]
	if !ok {
		return fmt.Errorf("command `%s` not yet implemented", cmd.Name)
	}
	err := r(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) registerCommand(name string, f func(*State, Command) error) {
	c.FuncFromCommand[name] = f
}

func middlewareLoggedIn(s *State, handler func(s *State, cmd Command) error) func(*State, Command) error {
	readConfig := config.Read()
	if readConfig.CurrentUsername == "" {
		log.Fatalln("error -> login first!")
	}
	return handler
}


func handlerFollowing(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, %s does not need any arguments\n", cmd.Name)
	}
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), string(s.Config_p.CurrentUsername))
	if err != nil {
		return err
	}
	fmt.Printf("List of feeds followed by user `%s`:\n", s.Config_p.CurrentUsername)
	for _, feed := range feeds {
		fmt.Printf("* %s -> %s\n", feed.FeedName, feed.FeedUrl)
	}
	return nil
}

func handlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error, %s needs additional arguments -> RSS URLs\n", cmd.Name)
	}
	_, err := s.Db.GetUser(context.Background(), s.Config_p.CurrentUsername)

	if err != nil {
		return fmt.Errorf("it seems user '%s' does not exist. is this user registered?", s.Config_p.CurrentUsername)
	}
	for _, url := range cmd.Args {
		err = _follow(s, url)
		if err != nil {
			return err
		}
	}
	return nil

}

func _follow(s *State, url string) error {
	fetchedUser, err := s.Db.GetUser(context.Background(), s.Config_p.CurrentUsername)
	if err != nil {
		return err
	}
	fetchedFeed, err := s.Db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    fetchedUser.ID,
		FeedID:    fetchedFeed.ID,
	}
	rows, err := s.Db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Println("Updated list of followed RSS feeds:")
	for _, row := range rows {
		fmt.Printf("User `%s` follows RSS feed `%s`\n", row.UserName, row.FeedName)
	}
	return nil
}

func handlerGetFeeds(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, %s does not need any arguments\n", cmd.Name)
	}
	_, err := s.Db.GetUser(context.Background(), s.Config_p.CurrentUsername)

	if err != nil {
		return fmt.Errorf("it seems user '%s' does not exist. is this user registered?", s.Config_p.CurrentUsername)
	}

	allFeeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	log.Println("Successfully got the list of feeds in the database:")
	for _, feed := range allFeeds {
		user, err := s.Db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to retrieve user from db: %v", err)
		}
		log.Printf("* %s -> %s, added by %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil

}

func handlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("error, %s needs two arguments -> a feed title and a feed URL\n", cmd.Name)
	}
	userInDb, err := s.Db.GetUser(context.Background(), s.Config_p.CurrentUsername)
	if err != nil {
    		return err
	}
	feedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    userInDb.ID,
	}

	_, err = s.Db.AddFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	err = _follow(s, cmd.Args[1])
	if err != nil {
		return err
	}
	return nil
}

func handlerAggregator(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, %s does not need any arguments\n", cmd.Name)
	}
	wagslane := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), wagslane)
	if err != nil {
		return err
	}
	log.Println(feed.Channel.Title)
	log.Println(feed.Channel.Link)
	log.Println(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		log.Println(item.Title)
		log.Println(item.Link)
		log.Println(item.Description)
		log.Println(item.PubDate)
	}
	return nil
}

func handlerGetUsers(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, %s does not need any arguments\n", cmd.Name)
	}
	allUsers, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	log.Println("Successfully got the list of users in the database:")
	for _, user := range allUsers {
		if s.Config_p.CurrentUsername == user.Name {
			log.Printf("* %s (current)\n", user.Name)
		} else {
			log.Printf("* %s\n", user.Name)
		}
	}
	return nil

}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 || cmd.Args == nil {
		return fmt.Errorf("error, %s needs additional arguments -> a name\n", cmd.Name)
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("error, %s only needs one argument -> a name\n", cmd.Name)
	}

	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("User '%s' is not registered. error: %v", cmd.Args[0], err)
	}
	(*s).Config_p.CurrentUsername = cmd.Args[0]
	(*s).Config_p.Write()
	log.Printf("User '%s' is logged in\n", cmd.Args[0])
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 || cmd.Args == nil {
		return fmt.Errorf("error, %s needs additional arguments\n -> a name", cmd.Name)
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("error, %s only needs one argument -> a name\n", cmd.Name)
	}

	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		return fmt.Errorf("error, user '%s' already exists.\n", cmd.Args[0])
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	_, err = s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	(*s).Config_p.CurrentUsername = cmd.Args[0]
	(*s).Config_p.Write()
	log.Printf("User '%s' is registered\n", cmd.Args[0])
	return nil
}

func handlerReset(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, %s does not need any arguments\n", cmd.Name)
	}
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	log.Println("Successfully reset the list of users in the database.")
	return nil
}
