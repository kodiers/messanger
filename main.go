package main

import (
	"messanger/libs/infrastructure/configuration"
	db2 "messanger/libs/infrastructure/database"
)

func main() {
	config := configuration.InitConfig()
	db := db2.ConnectToDb(config.GetDBConnectionString())
	db2.RunMigrations("migrations", db)
}
