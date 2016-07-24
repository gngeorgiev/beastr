package main

import (
	"log"

	"beatster-server/providers"

	"github.com/gin-gonic/gin"

	"sync"

	"net/http"

	"strings"
)

const (
	PARAM_QUERY    = "q"
	PARAM_ID       = "id"
	PARAM_PROVIDER = "provider"
)

func main() {
	r := gin.Default()

	player := r.Group("/player")
	{
		player.GET("/resolve", func(c *gin.Context) {
			id := c.Query(PARAM_ID)
			provider := strings.ToLower(c.Query(PARAM_PROVIDER))

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
			q := c.Query(PARAM_QUERY)
			registeredProviders := providers.GetProviders()
			result := map[string]interface{}{}

			wg := sync.WaitGroup{}
			wg.Add(len(registeredProviders))
			for _, p := range registeredProviders {
				go func(p providers.Provider) {
					res, err := p.Search(q)
					if err != nil {
						result[p.GetName()] = err
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

	log.Fatal(r.Run(":8080"))
}
