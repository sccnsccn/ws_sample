package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type StringJsonReader struct {
	Data []string
}

func (r *StringJsonReader) LoadData(fileNmae string) {
	jsonFile, err := os.Open(fileNmae)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal([]byte(byteValue), &r.Data)
}

// func main() {
// 	r := JsonReader{}
// 	r.loadData("valid_tokens.json")
// 	log.Println(r.validValue)
// }
