package controllers

import (
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"time"

	"log"

	"github.com/gin-gonic/gin"
)

const (
	AutocompleteUrl = "http://suggestqueries.google.com/complete/search?client=firefox&ds=yt&q=%s"
)

type autocompleteController struct {
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

func (a *autocompleteController) Autocomplete() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query(ParamQuery)
		autocompleteUrl := fmt.Sprintf(AutocompleteUrl, query)
		resp, err := http.Get(autocompleteUrl)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var autocompleteData []interface{}
		jsonErr := json.Unmarshal(body, &autocompleteData)
		if jsonErr != nil {
			log.Println(jsonErr)
			autocompleteData = make([]interface{}, 2)
			autocompleteData[0] = query
			autocompleteData[1] = make([]interface{}, 0)
		}

		cacheKey := a.GetCompleteCacheKey(c.Query(ParamQuery))
		cacheData(cacheKey, autocompleteData, time.Duration(24)*time.Hour)

		c.JSON(http.StatusOK, autocompleteData)
	}
}
