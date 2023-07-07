package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Token        string `json:"token"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Path         string `json:"path"`
	CustomHeader string `json:"customHeader"`
	Timeout      int    `json:"timeout"`
}

func (conf *Config) GetHost() string {
	return conf.Host + ":" + strconv.Itoa(conf.Port)
}

func (conf *Config) GetSecondsTimeout() time.Duration {
	return time.Duration(conf.Timeout) * time.Second
}

func (conf *Config) GetMillisecondsTimeout() time.Duration {
	return time.Duration(conf.Timeout) * time.Millisecond
}

func ReadConfig(fileNmae string) *Config {
	jsonFile, err := os.Open(fileNmae)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}

	var config Config

	json.Unmarshal([]byte(byteValue), &config)
	return &config
}

// func main() {
// 	config := readConfig("config.json")
// 	log.Println(config.CustomHeader)
// 	log.Println(config.Token)
// 	log.Println(config.Host)
// 	log.Println(config.Port)
// 	log.Println(config.CustomHeader)
// 	log.Println(config.getHost())
// }
