package providers

import (
	"net/http"

	"beatster-server/models"

	"fmt"

	"sync"

	"errors"

	"github.com/otium/ytdl"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

const (
	API_KEY = "AIzaSyDAfuetiKu8_xPk7TDmO5NlnYdkMoip8Tg"
)

type YouTubeProvider struct {
	provider

	service *youtube.Service
}

func (y *YouTubeProvider) GetService() *youtube.Service {
	return y.service
}

func (y *YouTubeProvider) getThumbnailUrl(t *youtube.ThumbnailDetails) string {
	if t.Maxres != nil {
		return t.Maxres.Url
	} else if t.High != nil {
		return t.High.Url
	} else if t.Medium != nil {
		return t.Medium.Url
	} else if t.Standard != nil {
		return t.Standard.Url
	} else if t.Default != nil {
		return t.Default.Url
	}

	//TODO: some default
	return ""
}

func (y *YouTubeProvider) getSpecificResults(kind string, items []*youtube.SearchResult) []*youtube.SearchResult {
	filteredItems := items[:0]
	for _, item := range items {
		if item.Id.Kind == kind {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

func (y *YouTubeProvider) Search(q string) ([]models.Track, error) {
	call := y.service.Search.List("id,snippet").Q(q).MaxResults(25)
	r, err := call.Do()
	if err != nil {
		return nil, err
	}

	videos := y.getSpecificResults("youtube#video", r.Items)
	results := make([]models.Track, len(videos))
	for i, item := range videos {
		if item.Id.Kind != "youtube#video" {
			continue
		}

		track := &models.Track{
			Id:        item.Id.VideoId,
			Provider:  y.GetName(),
			Thumbnail: y.getThumbnailUrl(item.Snippet.Thumbnails),
			Title:     item.Snippet.Title,
		}

		if i < len(videos)-1 {
			track.Next = videos[i+1].Id.VideoId
		} else {
			track.Next = videos[0].Id.VideoId
		}

		results[i] = *track
	}

	return results, nil
}

func (y *YouTubeProvider) GetUrlFromId(id string) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
}

func (y *YouTubeProvider) getVideoInfo(id string) (video *youtube.Video, err error) {
	call := y.service.Videos.List("snippet").Id(id).MaxResults(1)
	r, callError := call.Do()
	if callError != nil {
		err = callError
	}

	if r != nil && len(r.Items) > 0 {
		video = r.Items[0]
	}

	return
}

func (y *YouTubeProvider) getStreamUrl(id string) (string, error) {
	url := y.GetUrlFromId(id)
	info, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return "", err
	}

	var format ytdl.Format
	mp4Formats := info.Formats.
		Filter(ytdl.FormatResolutionKey, []interface{}{""}).
		Filter(ytdl.FormatVideoEncodingKey, []interface{}{""}).
		Filter(ytdl.FormatExtensionKey, []interface{}{"mp4"}).
		Best(ytdl.FormatAudioBitrateKey)

	if len(mp4Formats) > 0 {
		format = mp4Formats[0]
	} else {
		format = info.Formats.Best(ytdl.FormatAudioBitrateKey)[0]
	}

	downloadUrl, err := info.GetDownloadURL(format)
	if err != nil {
		return "", err
	}

	return downloadUrl.String(), nil
}

func (y *YouTubeProvider) getNextVideo(id string) (string, error) {
	res, err := y.GetService().Search.List("id").Type("video").RelatedToVideoId(id).Do()
	if err != nil {
		return "", err
	}

	for _, item := range res.Items {
		if item.Id.VideoId != id {
			return item.Id.VideoId, nil
		}
	}

	return res.Items[0].Id.VideoId, nil //we cannot find the next video, but lets still play one
}

func (y *YouTubeProvider) Resolve(id string) (models.Track, error) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	errs := make([]error, 0)

	//TODO: cache these calls
	var video *youtube.Video
	go func() {
		defer wg.Done()

		youtubeVideo, err := y.getVideoInfo(id)
		if err != nil {
			errs = append(errs, err)
		} else {
			video = youtubeVideo
		}
	}()

	var streamUrl string
	go func() {
		defer wg.Done()

		url, err := y.getStreamUrl(id)
		if err != nil {
			errs = append(errs, err)
		} else {
			streamUrl = url
		}
	}()

	var nextVideo string
	go func() {
		defer wg.Done()

		next, err := y.getNextVideo(id)
		if err != nil {
			errs = append(errs, err)
		} else {
			nextVideo = next
		}
	}()

	wg.Wait()
	if len(errs) > 0 {
		return models.Track{}, errors.New(fmt.Sprintf("%s", errs))
	}

	return models.Track{
		Id:        video.Id,
		Provider:  y.GetName(),
		StreamUrl: streamUrl,
		Thumbnail: y.getThumbnailUrl(video.Snippet.Thumbnails),
		Title:     video.Snippet.Title,
		Next:      nextVideo,
	}, nil
}

func init() {
	c := &http.Client{
		Transport: &transport.APIKey{Key: API_KEY},
	}

	s, err := youtube.New(c)
	if err != nil {
		panic(err)
	}

	registerProvider(&YouTubeProvider{
		provider: provider{
			domain: "youtube.com",
			name:   "YouTube",
		},
		service: s,
	})
}
