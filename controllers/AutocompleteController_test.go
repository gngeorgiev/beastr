package controllers

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestAutocompleteController(t *testing.T) {
	g := Goblin(t)

	autocompleteController := AutocompleteController
	g.Describe("Autocomplete", func() {
		g.Describe("Complete", func() {
			g.It("Should autocomplete hardwe to hardwell", func() {
				res, err := autocompleteController.autocomplete("hardwe")
				g.Assert(err == nil).IsTrue("Error during autocomplete")
				g.Assert(res[0].(string) == "hardwe").IsTrue("Wrong autocomplete first result, should be hardwe")
				g.Assert(len(res[1].([]interface{})) > 0).IsTrue("Incorrect autocomplete count")
			})

			g.It("Should return empty results when no query is passed", func() {
				res, err := autocompleteController.autocomplete("")
				g.Assert(err == nil).IsTrue("Error during autocomplete")
				g.Assert(res[0].(string) == "").IsTrue("Wrong autocomplete first result, should be empty")
				g.Assert(len(res[1].([]interface{})) == 0).IsTrue("Wrong autocomplete results, should be empty")
			})
		})
	})
}
