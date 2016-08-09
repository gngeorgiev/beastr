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
	baseController
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
			player.sendJson(c, http.StatusOK, data)
			return
		}

		id := c.Query(ParamId)
		provider := c.Query(ParamProvider)
		result, err := player.resolve(id, provider)
		if err != nil {
			player.sendError(c, http.StatusBadRequest, err)
			return
		}

		cacheKey := player.GetResolveCacheKey(id, provider)
		cacheData(cacheKey, result, time.Duration(1)*time.Hour)

		player.sendJson(c, http.StatusOK, result)
	}
}

func (p *playerController) search(q string) (result map[string]interface{}, hasErrors bool) {
	result = make(map[string]interface{})
	registeredProviders := providers.GetProviders()

	wg := sync.WaitGroup{}
	wg.Add(len(registeredProviders))
	for _, p := range registeredProviders {
		go func(p providers.Provider) {
			res, err := p.Search(q)
			if err != nil {
				result[p.GetName()] = err.Error()
				hasErrors = true
			} else {
				result[p.GetName()] = res
			}

			wg.Done()
		}(p)
	}

	wg.Wait()

	return
}

func (p *playerController) SearchRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, has := c.Get("cache")
		if has {
			p.sendJson(c, http.StatusOK, data)
			return
		}

		q := c.Query(ParamQuery)
		result, hasErrors := p.search(q)
		if !hasErrors {
			cacheKey := p.GetSearchCacheKey(q)
			cacheData(cacheKey, result, time.Duration(24)*time.Hour)
		}

		p.sendJson(c, http.StatusOK, result)
	}
}
