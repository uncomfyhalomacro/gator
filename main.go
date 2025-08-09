package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/uncomfyhalomacro/gator/internal/cli"
	"github.com/uncomfyhalomacro/gator/internal/config"
	"github.com/uncomfyhalomacro/gator/internal/database"
	"log"
	"os"
)

func main() {
	readConfig := config.Read()
	db, err := sql.Open("postgres", readConfig.DbUrl)
	if err != nil {
		log.Fatalf("%v", err)
	}
	dbQueries := database.New(db)
	commands := cli.Initialise()
	if len(os.Args) < 2 {
		log.Fatalf("gator requires a subcommand")
	}
	newCommand := cli.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	newState := cli.State{
		Db:       dbQueries,
		Config_p: &readConfig,
	}
	err = commands.Run(&newState, newCommand)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
