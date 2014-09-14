package main

import (
	"encoding/json"
	"github.com/dr4ke616/go_cloudflare"
	"log"
	"os"
	"time"
)

type Configuration struct {
	Name       string
	Frequencey time.Duration
	Cloudflare Cloudflare `json:"Cloudflare"`
}

type Cloudflare struct {
	Email      string
	Token      string
	Domain     string
	RecordID   string
	SubDomain  string
	RecordType string
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

func update(cloudflare *Cloudflare) {
	log.Println("About to check status of IP")

	client, err := go_cloudflare.NewClient(cloudflare.Email, cloudflare.Token)
	if err != nil {
		log.Fatal("Problem with clouflare client: ", err)
		os.Exit(1)
	}

	records, err := client.RetrieveARecord(cloudflare.Domain, cloudflare.RecordID)
	if err != nil {
		log.Fatal("Problem with clouflare client: ", err)
		os.Exit(1)
	}

	log.Println("Records:", records)

	err = client.UpdateRecord(cloudflare.Domain, cloudflare.RecordID, &go_cloudflare.UpdateRecord{
		Content: "80.111.125.147",
		Type:    cloudflare.RecordType,
		Name:    cloudflare.SubDomain,
	})
	if err != nil {
		log.Fatal("Problem with clouflare client: ", err)
		os.Exit(1)
	}
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
		update(&config.Cloudflare)
		time.Sleep(config.Frequencey * time.Second)
	}
}
