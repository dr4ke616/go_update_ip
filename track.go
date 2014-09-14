package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type IPInfo struct {
	IP        string
	Latlong   string
	Country   string
	City      string
	UserAgent string
}

func getIp() (*http.Response, error) {

	url := "http://www.trackip.net/ip?json"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request when tracking current IP: %s", err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error creating request when tracking current IP: %s", err)
	}

	return resp, nil
}

func decodeResponse(resp *http.Response, out interface{}) error {

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil

}

func trackCurrentIP() (*IPInfo, error) {

	resp, err := getIp()
	if err != nil {
		return nil, err
	}

	ipInfo := new(IPInfo)
	err = decodeResponse(resp, &ipInfo)
	if err != nil {
		return nil, fmt.Errorf("Problem decoding IP response", err)
	}

	return ipInfo, nil
}
