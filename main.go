package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gngeorgiev/beatstr-server/controllers"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetDefault("redis_address", "localhost:6379")
	viper.BindEnv("redis_address")
}

func main() {
	initConfig()

	r := gin.Default()

	r.Use(controllers.GetMiddleware()...)

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

	log.Fatal(r.Run(":8085"))
}
