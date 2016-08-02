package main

import (
	"log"

	"beatster-server/providers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"sync"

	"net/http"

	"beatster-server/clients"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	ParamQuery    = "q"
	ParamId       = "id"
	ParamProvider = "provider"
)

func initConfig() {
	viper.SetDefault("redis_address", "localhost:6379")
	viper.BindEnv("redis_address")
}

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

func main() {
	initConfig()

	r := gin.Default()

	player := r.Group("/player")
	{
		player.Use(func(c *gin.Context) {
			var cacheKey string

			url := c.Request.URL.String()
			if strings.Contains(url, "/resolve") {
				cacheKey = fmt.Sprintf("%s%s", c.Query(ParamId), c.Query(ParamProvider))
			} else if strings.Contains(url, "/search") {
				cacheKey = c.Query(ParamQuery)
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
		})

		player.GET("/resolve", func(c *gin.Context) {
			data, has := c.Get("cache")
			if has {
				c.JSON(http.StatusOK, data)
				return
			}

			id := c.Query(ParamId)
			provider := strings.ToLower(c.Query(ParamProvider))

			registeredProviders := providers.GetProviders()
			for _, p := range registeredProviders {
				if strings.ToLower(p.GetName()) != provider {
					continue
				}

				result, err := p.Resolve(id)
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}

				cacheKey := fmt.Sprintf("%s%s", c.Query(ParamId), c.Query(ParamProvider))
				cacheData(cacheKey, result, time.Duration(1)*time.Hour)

				c.JSON(http.StatusOK, result)
				return
			}

			c.AbortWithStatus(http.StatusNotFound)
		})

		player.GET("/search", func(c *gin.Context) {
			data, has := c.Get("cache")
			if has {
				c.JSON(http.StatusOK, data)
				return
			}

			q := c.Query(ParamQuery)
			result := map[string]interface{}{}
			registeredProviders := providers.GetProviders()

			wg := sync.WaitGroup{}
			wg.Add(len(registeredProviders))
			for _, p := range registeredProviders {
				go func(p providers.Provider) {
					res, err := p.Search(q)
					if err != nil {
						result[p.GetName()] = err.Error()
					} else {
						result[p.GetName()] = res
					}

					wg.Done()
				}(p)
			}

			wg.Wait()

			cacheData(q, result, time.Duration(24)*time.Hour)

			c.JSON(http.StatusOK, result)
		})
	}

	log.Fatal(r.Run(":8085"))
}
