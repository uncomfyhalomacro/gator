package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func Read() Config {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error occured. error: %v", err)
	}
	configFilePath := filepath.Join(homedir, configFileName)
	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("cannot read config file. error: %v", err)
	}
	defer configFile.Close()
	var newConfig Config
	buf := make([]byte, 1024)
	var actualReadN int
	for {
		n, err := configFile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		actualReadN = n
	}

	errJ := json.Unmarshal(buf[:actualReadN], &newConfig)
	if errJ != nil {
		log.Fatalf("an error occured: %v", errJ)
	}

	return newConfig
}
