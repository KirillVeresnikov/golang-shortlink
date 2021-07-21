package app

import (
	"golang-shortlink/pkg/shortLinkService/infrastructure/jsonProvider"
)

type jsonProviderInterface interface {
	LoadPaths(src string)
	GetURL(shortURL string) string
	GetErr() error
}

type Service struct {
	jsonProvider jsonProviderInterface
	src          string
	Err          error
}

func Create(src string) *Service {
	service := Service{}
	service.jsonProvider = jsonProvider.Create()
	service.src = src
	service.jsonProvider.LoadPaths(service.src)
	return &service
}

func (s *Service) GetLongURL(shortURL string) string {
	longURL := s.jsonProvider.GetURL(shortURL)
	if err := s.jsonProvider.GetErr(); err != nil {
		s.Err = err
		return ""
	}
	s.Err = nil
	return longURL
}

func (s *Service) GetErr() error {
	return s.Err
}
