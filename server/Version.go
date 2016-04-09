package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Version struct {
	Version string `json:"version"`
}

func GetServerVersionFromMaster() (*Version, error) {
	var result *Version

	client := http.Client{}
	response, err := client.Get("https://raw.githubusercontent.com/mailslurper/mailslurper/master/version.json")

	if err != nil {
		return result, err
	}

	versionBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(versionBytes, &result); err != nil {
		return result, err
	}

	return result, nil
}
