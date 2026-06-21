package main

import (
	"log"
	"os"

	"github.com/rizalarfiyan/be-plant-factory/internal/bootstrap"
)

// Swagger
//
//	@title						BE Plant Factory API
//	@version					1.0
//	@description				This is an API documentation of BE Plant Factory.
//	@BasePath					/
//	@securityDefinitions.apikey	CookieAccessToken
//	@in							cookie
//	@name						access_token
//	@securityDefinitions.apikey	CookieRefreshToken
//	@in							cookie
//	@name						refresh_token
func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("application stopped with error: %v", err)
	}

	os.Exit(0)
}
