package main

import (
	"github.com/uncomfyhalomacro/gator/internal/cli"
	"github.com/uncomfyhalomacro/gator/internal/config"
	"log"
	"os"
)

func main() {
	readConfig := config.Read()
	commands := cli.Initialise()
	if len(os.Args) == 1 {
		log.Fatalf("gator requires a subcommand")
	}
	if len(os.Args) == 2 {
		log.Fatalf("gator subcommand requires additional arguments")
	}
	newCommand := cli.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	newState := cli.State{
		Config_p: &readConfig,
	}
	err := commands.Run(&newState, newCommand)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
