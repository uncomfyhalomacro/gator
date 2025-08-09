package cli

import (
	"fmt"
	"github.com/uncomfyhalomacro/gator/internal/config"
)

type State struct {
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
	c.register("login", handlerLogin)
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

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.FuncFromCommand[name] = f
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 || cmd.Args == nil {
		return fmt.Errorf("error, %s needs additional arguments\n", cmd.Name)
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("error, %s only needs one argument\n", cmd.Name)
	}

	state := *s
	config := *state.Config_p
	config.CurrentUsername = cmd.Args[0]
	config.Write()
	return nil
}
