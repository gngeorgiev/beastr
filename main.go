package main

import (
	"log"

	"fmt"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/gngeorgiev/beatster-server/clients"
	"github.com/gngeorgiev/beatster-server/controllers"
	"github.com/spf13/viper"
)

var version string
var hostname string

func initConfig() {
	viper.SetDefault("redis_address", "localhost:6379")
	viper.BindEnv("redis_address")
}

func initVariables() {
	if version == "" {
		version = "development"
	}

	osHostname, err := os.Hostname()
	if err != nil {
		hostname = err.Error()
	} else {
		hostname = osHostname
	}
}

func initLogging() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func initServices() {
	clients.StartRedisConnection()
}

func main() {
	initVariables()
	initLogging()
	initConfig()
	initServices()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		gin.Logger()(c)

		query := c.Request.URL.Query()

		fmt.Println(fmt.Sprintf(`Query: %s,`,
			query,
		))
	})

	mainController := controllers.MainController

	r.Use(mainController.GetMiddleware()...)
	r.GET("/status", mainController.StatusRouteHandler(version, hostname))

	playerController := controllers.PlayerController
	player := r.Group(playerController.GetPrefix())
	{
		player.Use(playerController.GetMiddleware()...)
		player.GET("/resolve", playerController.ResolveRouteHandler())
		player.GET("/search", playerController.SearchRouteHandler())
	}

	autocompleteController := controllers.AutocompleteController
	autocomplete := r.Group(autocompleteController.GetPrefix())
	{
		autocomplete.Use(autocompleteController.GetMiddleware()...)
		autocomplete.GET("/complete", autocompleteController.AutocompleteRouteHandler())
	}

	log.Println(fmt.Sprintf("Server version: \"%s\"", version))
	log.Println(fmt.Sprintf("Server hostname: \"%s\"", hostname))
	log.Println(fmt.Sprintf("Redis address: \"%s\"", viper.GetString("redis_address")))

	log.Fatal(r.Run(":8085"))
}
