package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"beatster-server/providers"
	"time"

	"sync"

	"github.com/gin-gonic/gin"
)

type playerController struct {
}

func newPlayerController() *playerController {
	return &playerController{}
}

func (p *playerController) GetPrefix() string {
	return "/player"
}

func (p *playerController) GetMiddleware() []gin.HandlerFunc {
	middleware := make([]gin.HandlerFunc, 0)
	return middleware
}

func (p *playerController) GetResolveCacheKey(id, provider string) string {
	return fmt.Sprintf("resolve_%s_%s", id, provider)
}

func (p *playerController) GetSearchCacheKey(query string) string {
	return fmt.Sprintf("search_%s", query)
}

func (player *playerController) Resolve() gin.HandlerFunc {
	return func(c *gin.Context) {
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

			cacheKey := player.GetResolveCacheKey(c.Query(ParamId), c.Query(ParamProvider))
			cacheData(cacheKey, result, time.Duration(1)*time.Hour)

			c.JSON(http.StatusOK, result)
			return
		}

		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (p *playerController) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		cacheKey := p.GetSearchCacheKey(q)
		cacheData(cacheKey, result, time.Duration(24)*time.Hour)

		c.JSON(http.StatusOK, result)
	}
}
