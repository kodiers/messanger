package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"messanger/libs/infrastructure/configuration"
	"messanger/libs/infrastructure/database"
	"messanger/users"
	"net/http"
)

func main() {
	// Preparing to launch
	database.RunMigrations("migrations", configuration.DB)

	// Add routes
	router := httprouter.New()
	router.POST("/users/register", users.Register)
	router.POST("/users/login", users.Login)

	log.Fatalln(http.ListenAndServe(":80", router))
}
