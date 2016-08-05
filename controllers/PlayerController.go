package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"time"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gngeorgiev/beatstr-server/models"
	"github.com/gngeorgiev/beatstr-server/providers"
	"github.com/go-errors/errors"
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

func (player *playerController) resolve(id, provider string) (models.Track, error) {
	provider = strings.ToLower(provider)
	registeredProviders := providers.GetProviders()
	for _, p := range registeredProviders {
		if strings.ToLower(p.GetName()) != provider {
			continue
		}

		return p.Resolve(id)
	}

	return models.Track{}, errors.New("Not found")
}

func (player *playerController) ResolveRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, has := c.Get("cache")
		if has {
			c.JSON(http.StatusOK, data)
			return
		}

		id := c.Query(ParamId)
		provider := c.Query(ParamProvider)
		result, err := player.resolve(id, provider)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		cacheKey := player.GetResolveCacheKey(id, provider)
		cacheData(cacheKey, result, time.Duration(1)*time.Hour)

		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (p *playerController) search(q string) map[string]interface{} {
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

	return result
}

func (p *playerController) SearchRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, has := c.Get("cache")
		if has {
			c.JSON(http.StatusOK, data)
			return
		}

		q := c.Query(ParamQuery)
		result := p.search(q)

		cacheKey := p.GetSearchCacheKey(q)
		cacheData(cacheKey, result, time.Duration(24)*time.Hour)

		c.JSON(http.StatusOK, result)
	}
}
