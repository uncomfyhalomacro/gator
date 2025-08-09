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
		log.Fatalf("cannot open config file. error: %v", err)
	}
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
	err = configFile.Close()
	if err != nil {
		log.Fatalf("error occured when closing file. error: %v", err)
	}

	return newConfig
}

func (c *Config) Write() error {
	// Marshalling back to json bytes
	buf, err := json.Marshal(*c)
	if err != nil {
		return fmt.Errorf("error occured while reading JSON config: %v", err)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error occured. error: %v", err)
	}
	configFilePath := filepath.Join(homedir, configFileName)
	configFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModeAppend) // Overwrites instead of appends
	if err != nil {
		log.Fatalf("cannot open config file. error: %v", err)
	}
	n, err := configFile.Write(buf)
	if err != nil {
		log.Fatalf("error occured while writing into file. error: %v", err)
	}
	log.Printf("wrote %d bytes to config file\n", n)
	err = configFile.Close()
	if err != nil {
		log.Fatalf("error occured when closing file. error: %v", err)
	}
	return nil
}
