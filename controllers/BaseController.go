package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"strings"

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

func GetMiddleware() []gin.HandlerFunc {
	middleware := make([]gin.HandlerFunc, 1)

	middleware[0] = func(c *gin.Context) {
		var cacheKey string

		url := c.Request.URL.String()
		if strings.Contains(url, "/resolve") {
			cacheKey = PlayerController.GetResolveCacheKey(c.Query(ParamId), c.Query(ParamProvider))
		} else if strings.Contains(url, "/search") {
			cacheKey = PlayerController.GetSearchCacheKey(c.Query(ParamQuery))
		} else if strings.Contains(url, "/complete") {
			cacheKey = AutocompleteController.GetCompleteCacheKey(c.Query(ParamQuery))
		}

		var result interface{}
		redisClient := clients.GetRedisClient()
		cachedData, err := redisClient.Get(cacheKey).Result()
		if err == nil && cachedData != "" {
			jsonErr := json.Unmarshal([]byte(cachedData), &result)
			if jsonErr != nil {
				c.AbortWithError(http.StatusInternalServerError, jsonErr)
				return
			}

			c.Set("cache", result)
		} else if err != nil {
			log.Println(err)
		}
	}

	return middleware
}
