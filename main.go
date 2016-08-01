package main

import (
	"log"

	"beatster-server/providers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"sync"

	"net/http"

	"strings"
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

func main() {
	initConfig()

	r := gin.Default()

	player := r.Group("/player")
	{
		player.GET("/resolve", func(c *gin.Context) {
			id := c.Query(ParamId)
			provider := strings.ToLower(c.Query(ParamProvider))

			registeredProviders := providers.GetProviders()
			for _, p := range registeredProviders {
				if strings.ToLower(p.GetName()) != provider {
					continue
				}

				t, err := p.Resolve(id)
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}

				c.JSON(http.StatusOK, t)
				return
			}

			c.AbortWithStatus(http.StatusNotFound)
		})

		player.GET("/search", func(c *gin.Context) {
			q := c.Query(ParamQuery)
			registeredProviders := providers.GetProviders()
			result := map[string]interface{}{}

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

			c.JSON(http.StatusOK, result)
		})
	}

	log.Fatal(r.Run(":8085"))
}
