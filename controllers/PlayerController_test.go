package controllers

import (
	"testing"

	"fmt"

	. "github.com/franela/goblin"
	"github.com/gngeorgiev/beatstr-server/models"
)

func TestPlayerController(t *testing.T) {
	g := Goblin(t)

	checkTrack := func(track models.Track, id string, checkStreamUrl bool) {
		if id != "" {
			g.Assert(track.Id == id).IsTrue("Wrong id " + track.Id)
		}

		if checkStreamUrl {
			g.Assert(track.StreamUrl != "").IsTrue("No streamurl found")
		}

		g.Assert(track.Next != "").IsTrue("No next video found")
		g.Assert(track.Provider == "YouTube").IsTrue("Wrong provider " + track.Provider)
		g.Assert(track.Thumbnail != "").IsTrue("No thuumbnail")
		g.Assert(track.Title != "").IsTrue("No title")
	}

	playerController := PlayerController
	g.Describe("Player", func() {
		g.Describe("Resolve", func() {
			g.It("Resolve Don't let me down", func() {
				track, err := playerController.resolve("o0citpYDaVg", "youtube")
				g.Assert(err == nil).IsTrue("Error during resolve")
				checkTrack(track, "o0citpYDaVg", true)
				g.Assert(track.Title == "The Chainsmokers - Don't Let Me Down (Lyric) ft. Daya").IsTrue("Wrong title " + track.Title)
			})
		})

		g.Describe("Search", func() {
			g.It("Search chainsmokers", func() {
				results, _ := playerController.search("chainsmokers")
				youTube := results["YouTube"]
				if err, ok := youTube.(string); ok {
					g.Fail("There was an error while searching" + err)
				}

				youtubeResults, _ := youTube.([]models.Track)

				g.Assert(len(youtubeResults) > 0).IsTrue(fmt.Sprintf("Wrong results count %d", len(youtubeResults)))

				for _, track := range youtubeResults {
					checkTrack(track, "", false)
				}
			})
		})
	})
}
