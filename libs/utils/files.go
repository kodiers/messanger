package utils

import (
	"log"
	"os"
)

func ReadConfig(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		log.Fatalln("Could not read file at: ", filePath)
		return nil, err
	}
	var data []byte
	_, err = file.Read(data)
	if err != nil {
		log.Fatalln("Could not read data from file at: ", filePath)
		return nil, err
	}
	return data, nil
}

func Contains(sl []string, e string) bool {
	for _, s := range sl {
		if s == e {
			return true
		}
	}
	return false
}
