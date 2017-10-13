package ari

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ariURI string = "http://127.0.0.1:8088"
)

// Endpoints define an endpoints
type Endpoints struct {
	ChannelIds string `json:"channel_ids"`
	Resource   string `json:"resource"`
	State      string `json:"state"`
	Technology string `json:"technology"`
}

// ModuleReload reload a module
func ModuleReload(module string) error {
	url := ariURI + "/asterisk/modules/" + module

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, reqErr := http.NewRequest(http.MethodPut, url, nil)
	if reqErr != nil {
		return reqErr
	}

	req.Header.Set("User-Agent", "Golang Echo")

	res, doErr := client.Do(req)
	if doErr != nil {
		return doErr
	}

	if res.StatusCode == 404 {
		return errors.New("Not found")
	} else if res.StatusCode == 409 {
		return errors.New("Reload failed")
	}

	return nil
}

// ListEndpoints available endoints for a given endpoint technology
func ListEndpoints(technology string) (list []Endpoints, err error) {
	url := ariURI + "/endpoints/" + technology

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("User-Agent", "Golang Echo")

	res, doErr := client.Do(req)
	if doErr != nil {
		return nil, doErr
	}

	if res.StatusCode == 404 {
		return nil, errors.New("Not found")
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	jsonErr := json.Unmarshal(body, &list)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return list, nil
}
