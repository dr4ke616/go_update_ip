package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Configuration struct {
	Name       string
	Frequencey time.Duration
}

func load_configuration(config ...string) (c Configuration, err_ error) {
	var err error
	var file *os.File
	var config_file = "config/application.json"

	if len(config) > 0 {
		config_file = config[0]
	}

	if file, err = os.Open(config_file); err != nil {
		log.Println("Failed to open config file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	if err = decoder.Decode(&configuration); err != nil {
		log.Println("Failed to decode JSON config file:", err)
		return
	}

	return configuration, nil
}

func update() {
	log.Println("About to check status of IP")

}

func main() {
	log.Println("Dynamic IP updater started")

	config, err := load_configuration()
	if err != nil {
		log.Fatal("Failed to load config file. Abort!!")
		os.Exit(1)
	}

	log.Println("Starting", config.Name)
	for {
		update()
		time.Sleep(config.Frequencey * time.Second)
	}
}
