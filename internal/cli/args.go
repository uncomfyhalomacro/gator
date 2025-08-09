package cli

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/uncomfyhalomacro/gator/internal/config"
	"github.com/uncomfyhalomacro/gator/internal/database"
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

func Initialise() Commands {
	mapping := make(map[string]func(*State, Command) error)
	c := Commands{}
	c.FuncFromCommand = mapping
	c.registerCommand("login", handlerLogin)
	c.registerCommand("register", handlerRegister)
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

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 || cmd.Args == nil {
		return fmt.Errorf("error, %s needs additional arguments -> a name\n", cmd.Name)
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("error, %s only needs one argument -> a name\n", cmd.Name)
	}

	state := *s
	_, err := state.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	state.Config_p.CurrentUsername = cmd.Args[0]
	state.Config_p.Write()
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 || cmd.Args == nil {
		return fmt.Errorf("error, %s needs additional arguments\n -> a name", cmd.Name)
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("error, %s only needs one argument -> a name\n", cmd.Name)
	}

	state := *s
	_, err := state.Db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		return fmt.Errorf("error, user '%s' already exists.\n", cmd.Args[0])
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	_, err = state.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	state.Config_p.CurrentUsername = cmd.Args[0]
	state.Config_p.Write()
	return nil
}
