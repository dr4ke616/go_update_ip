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
	Name       string
	RecordType string
}

func LoadConfiguration(config ...string) (c Configuration, err_ error) {
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

func Update(cloudflare *Cloudflare) {
	log.Println("About to check status of IP")

	ipInfo, err := trackCurrentIP()
	if err != nil {
		log.Fatal("Problem getting current IP Address:", err)
		os.Exit(1)
	}

	log.Println("Your current IP info:")
	log.Println("IP: ", ipInfo.IP)
	log.Println("Latitude/longitude: ", ipInfo.Latlong)
	log.Println("Country: ", ipInfo.Country)
	log.Println("City: ", ipInfo.City)

	client, err := go_cloudflare.NewClient(cloudflare.Email, cloudflare.Token)
	if err != nil {
		log.Fatal("Problem with clouflare client: ", err)
		os.Exit(1)
	}

	record, err := client.RetrieveARecord(cloudflare.Domain, cloudflare.RecordID)
	if err != nil {
		log.Fatal("Problem with clouflare client: ", err)
		os.Exit(1)
	}

	if record.Value == ipInfo.IP {
		log.Println("Ip address match. No need to update")
		return
	}

	err = client.UpdateRecord(cloudflare.Domain, cloudflare.RecordID, &go_cloudflare.UpdateRecord{
		Content: ipInfo.IP,
		Type:    cloudflare.RecordType,
		Name:    cloudflare.Name,
	})
	if err != nil {
		log.Fatal("Problem with clouflare client: ", err)
		os.Exit(1)
	}

	log.Println("Ip address were different, so I updated it.")
}

func main() {
	log.Println("Dynamic IP updater started")

	config, err := LoadConfiguration()
	if err != nil {
		log.Fatal("Failed to load config file. Abort!!")
		os.Exit(1)
	}

	log.Println("Starting", config.Name)
	for {
		Update(&config.Cloudflare)
		time.Sleep(config.Frequencey * time.Second)
	}
}
