package utils

import (
	"io/ioutil"
	"log"
)

func ReadConfig(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Could not read data from file at: ", filePath)
		return nil, err
	}
	return data, nil
}
