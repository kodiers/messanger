package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"messanger/libs/infrastructure/configuration"
	"messanger/libs/infrastructure/database"
	"messanger/messages"
	"messanger/users"
	"net/http"
)

func main() {
	// Preparing to launch
	database.RunMigrations("migrations", configuration.DB)

	// Add routes
	router := httprouter.New()
	// User routes
	router.POST("/users/register", users.Register)
	router.POST("/users/login", users.Login)
	router.GET("/users/list", users.AuthenticationMiddleware(users.UsersList))

	// messages routes
	router.POST("/messages/create", users.AuthenticationMiddleware(messages.CreateMessage))

	log.Fatalln(http.ListenAndServe(":80", router))
}
