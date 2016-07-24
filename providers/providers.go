package providers

import (
	"beatster-server/models"
	"errors"
)

type Provider interface {
	GetDomain() string
	GetName() string
	Search(query string) ([]models.Track, error)
	Resolve(id string) (models.Track, error)
	GetUrlFromId(id string) string
}

type provider struct {
	domain, name string
}

func (p *provider) GetDomain() string {
	return p.domain
}

func (p *provider) GetName() string {
	return p.name
}

func (p *provider) Search(q string) ([]models.Track, error) {
	return nil, errors.New("Override Search")
}

func (p *provider) Resolve(id string) (models.Track, error) {
	return models.Track{}, errors.New("Override Resolve")
}

func (p *provider) GetStringFromId(id string) string {
	return ""
}

var registeredProviders = make([]Provider, 0)

func registerProvider(p Provider) {
	registeredProviders = append(registeredProviders, p)
}

func GetProviders() []Provider {
	return registeredProviders
}
