package main

import (
	_ "github.com/lib/pq"
	"github.com/uncomfyhalomacro/gator/internal/cli"
	"log"
	"os"
)

func listCommands(commands cli.Commands) string {
	var commandListString string
	for key, _ := range commands.FuncFromCommand {
		commandListString = commandListString + "* " + key + "\n"
	}
	return commandListString
}

func main() {
	commands := cli.Initialise()
	if len(os.Args) < 2 {
		log.Fatalf("gator requires a subcommand. available subcommands:\n%s", listCommands(commands))
	}
	newCommand := cli.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	err := commands.Run(newCommand)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
