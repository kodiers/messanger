package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"messanger/libs/utils"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"database"`
}

type MigrationsConfig struct {
	FolderPath string `yaml:"folder_path"`
}

type MetaConfig struct {
	Migrations MigrationsConfig `yaml:"migrations"`
}

type Config struct {
	DB   DBConfig   `yaml:"db"`
	Meta MetaConfig `yaml:"meta"`
}

func InitConfig() Config {
	configData := Config{}
	configData.getConfig()
	return configData
}

func (c *Config) getConfig() {
	data, err := utils.ReadConfig("config/config.yml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		log.Fatalln("Error reading configuration!")
	}
}

func (c Config) GetDBConnectionString() string {
	result := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", c.DB.User, c.DB.Password,
		c.DB.Host, c.DB.Port, c.DB.DBName)
	return result
}

func (c Config) GetPathToMigrationsFolder() string {
	return c.Meta.Migrations.FolderPath
}
