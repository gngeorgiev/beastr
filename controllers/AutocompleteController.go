package controllers

import (
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"time"

	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	AutocompleteUrl = "http://suggestqueries.google.com/complete/search?client=firefox&ds=yt&q=%s"
)

type autocompleteController struct {
	baseController
}

func newAutocompleteController() *autocompleteController {
	return &autocompleteController{}
}

func (a *autocompleteController) GetPrefix() string {
	return "/autocomplete"
}

func (a *autocompleteController) GetMiddleware() []gin.HandlerFunc {
	middleware := make([]gin.HandlerFunc, 0)
	return middleware
}

func (a *autocompleteController) GetCompleteCacheKey(query string) string {
	return fmt.Sprintf("complete_%s", query)
}

func (a *autocompleteController) autocomplete(query string) ([]interface{}, error) {
	if query == "" {
		defaultResult := make([]interface{}, 2)
		defaultResult[0] = query
		defaultResult[1] = make([]interface{}, 0)
		return defaultResult, nil
	}

	autocompleteUrl := fmt.Sprintf(AutocompleteUrl, url.QueryEscape(query))
	resp, err := http.Get(autocompleteUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var autocompleteData []interface{}
	jsonErr := json.Unmarshal(body, &autocompleteData)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return autocompleteData, nil
}

func (a *autocompleteController) AutocompleteRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, has := c.Get("cache")
		if has {
			a.sendJson(c, http.StatusOK, data)
			return
		}

		query := c.Query(ParamQuery)
		autocompleteData, err := a.autocomplete(query)
		if err != nil {
			a.sendError(c, http.StatusBadRequest, err)
			return
		}

		cacheKey := a.GetCompleteCacheKey(query)
		cacheData(cacheKey, autocompleteData, time.Duration(24)*time.Hour)

		a.sendJson(c, http.StatusOK, autocompleteData)
	}
}
