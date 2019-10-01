package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"database"`
}

type Config struct {
	DB DBConfig `yaml:"db"`
}

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

func GetDBConnectionString() (string, error) {
	configData := Config{}
	data, err := ReadConfig("config/config.yml")
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(data, &configData)
	if err != nil {
		log.Fatalln("Error reading configuration!")
		return "", err
	}
	result := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", configData.DB.User, configData.DB.Password,
		configData.DB.Host, configData.DB.Port, configData.DB.DBName)
	return result, nil
}
