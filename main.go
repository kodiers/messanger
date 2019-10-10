package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"messanger/libs/infrastructure/configuration"
	db2 "messanger/libs/infrastructure/database"
	"messanger/users"
	"net/http"
)

func main() {
	// Preparing to launch
	config := configuration.InitConfig()
	db := db2.ConnectToDb(config.GetDBConnectionString())
	db2.RunMigrations("migrations", db)

	// Add routes
	router := httprouter.New()
	router.GET("/users/register", users.Register)

	log.Fatalln(http.ListenAndServe(":80", router))
}
