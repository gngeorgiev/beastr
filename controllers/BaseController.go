package controllers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gngeorgiev/beatstr-server/clients"
)

const (
	ParamQuery    = "q"
	ParamId       = "id"
	ParamProvider = "provider"
)

var (
	PlayerController       = newPlayerController()
	AutocompleteController = newAutocompleteController()
	MainController         = newMainController()
)

func cacheData(key string, data interface{}, duration time.Duration) {
	redisClient := clients.GetRedisClient()
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	} else {
		err = redisClient.Set(key, string(b), duration).Err()
		if err != nil {
			log.Println(err)
		}
	}
}

type BaseController interface {
	GetPrefix() string
	GetMiddleware() []gin.HandlerFunc
}
