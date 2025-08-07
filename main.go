package main

import (
	"encoding/json"
	"fmt"
	"github.com/uncomfyhalomacro/gator/internal/config"
	"os"
)

func main() {
	gatorConfig := config.Read()
	// Marshalling back to json bytes
	buf, err := json.Marshal(gatorConfig)
	if err != nil {
		fmt.Println("an error occured: %v", err)
		os.Exit(2)
	}
	fmt.Println(string(buf))

}
