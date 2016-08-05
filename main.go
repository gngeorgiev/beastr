package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gngeorgiev/beatstr-server/controllers"
	"github.com/spf13/viper"
	"net/http"
	"fmt"
)

func initConfig() {
	viper.SetDefault("redis_address", "localhost:6379")
	viper.BindEnv("redis_address")
}

var version string

func main() {
	if version == "" {
		version = "development"
	}

	initConfig()

	r := gin.Default()

	r.Use(controllers.GetMiddleware()...)

	r.GET("/status", func (c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"version": version,
		})
	})

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

	log.Fatal(r.Run(":8085"))
}
