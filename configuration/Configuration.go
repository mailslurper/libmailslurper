// Copyright 2013-2014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package configuration

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/adampresley/golangdb"
)

/*
The Configuration structure represents a JSON
configuration file with settings for how to bind
servers and connect to databases.
*/
type Configuration struct {
	WWWAddress     string `json:"wwwAddress"`
	WWWPort        int    `json:"wwwPort"`
	ServiceAddress string `json:"serviceAddress"`
	ServicePort    int    `json:"servicePort"`
	SmtpAddress    string `json:"smtpAddress"`
	SmtpPort       int    `json:"smtpPort"`
	DBEngine       string `json:"dbEngine"`
	DBHost         string `json:"dbHost"`
	DBPort         int    `json:"dbPort"`
	DBDatabase     string `json:"dbDatabase"`
	DBUserName     string `json:"dbUserName"`
	DBPassword     string `json:"dbPassword"`
	MaxWorkers     int    `json:"maxWorkers"`
}

/*
Returns a pointer to a DatabaseConnection structure with data
pulled from a Configuration structure.
*/
func (this *Configuration) GetDatabaseConfiguration() *golangdb.DatabaseConnection {
	return &golangdb.DatabaseConnection{
		Engine:   golangdb.GetDatabaseEngineFromName(this.DBEngine),
		Address:  this.DBHost,
		Port:     this.DBPort,
		Database: this.DBDatabase,
		UserName: this.DBUserName,
		Password: this.DBPassword,
	}
}

/*
Returns a full address and port for the MailSlurper service
application.
*/
func (this *Configuration) GetFullServiceAppAddress() string {
	return fmt.Sprintf("%s:%d", this.ServiceAddress, this.ServicePort)
}

/*
Returns a full address and port for the MailSlurper SMTP
server.
*/
func (this *Configuration) GetFullSmtpBindingAddress() string {
	return fmt.Sprintf("%s:%d", this.SmtpAddress, this.SmtpPort)
}

/*
Returns a full address and port for the Web application.
*/
func (this *Configuration) GetFullWwwBindingAddress() string {
	return fmt.Sprintf("%s:%d", this.WWWAddress, this.WWWPort)
}

/*
Reads data from a Reader into a new Configuration structure.
*/
func LoadConfiguration(reader io.Reader) (*Configuration, error) {
	var err error
	var contents bytes.Buffer
	var buffer = make([]byte, 4096)
	var bytesRead int

	result := &Configuration{}
	bufferedReader := bufio.NewReader(reader)

	for {
		bytesRead, err = bufferedReader.Read(buffer)
		if err != nil && err != io.EOF {
			return result, err
		}

		if bytesRead == 0 {
			break
		}

		if _, err := contents.Write(buffer[:bytesRead]); err != nil {
			return result, err
		}
	}

	err = json.Unmarshal(contents.Bytes(), result)
	if err != nil {
		return result, err
	}

	return result, nil
}

/*
Reads data from a file into a Configuration object. Makes use of
LoadConfiguration().
*/
func LoadConfigurationFromFile(fileName string) (*Configuration, error) {
	result := &Configuration{}

	configFileHandle, err := os.Open(fileName)
	if err != nil {
		return result, err
	}

	result, err = LoadConfiguration(configFileHandle)
	if err != nil {
		return result, err
	}

	return result, nil
}

/*
Saves the current state of a Configuration structure
into a JSON file.
*/
func (this *Configuration) SaveConfiguration(configFile string) error {
	json, err := json.Marshal(this)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFile, json, 0644)
	if err != nil {
		return err
	}

	return nil
}
