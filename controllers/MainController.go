package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/gngeorgiev/beatstr-server/clients"
	"gopkg.in/redis.v4"
)

type mainController struct {
	baseController
}

func newMainController() *mainController {
	return &mainController{}
}

func (m *mainController) GetPrefix() string {
	return ""
}

func (m *mainController) GetMiddleware() []gin.HandlerFunc {
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
		} else {
			c.Next()
			return
		}

		var result interface{}
		redisClient := clients.GetRedisClient()
		cachedData, err := redisClient.Get(cacheKey).Result()
		if cachedData != "" {
			jsonErr := json.Unmarshal([]byte(cachedData), &result)
			if jsonErr != nil {
				m.sendError(c, http.StatusInternalServerError, jsonErr)
				return
			}

			c.Set("cache", result)
		} else if err != nil && err != redis.Nil {
			log.Println(err)
		}
	}

	return middleware
}

func (m *mainController) status(version, hostname string) map[string]interface{} {
	timestamp := time.Now().Format(time.StampMilli)

	return map[string]interface{}{
		"version":   version,
		"hostname":  hostname,
		"timestamp": timestamp,
	}
}

func (m *mainController) StatusRouteHandler(version, hostname string) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := m.status(version, hostname)
		c.JSON(http.StatusOK, status)
	}
}
