package jsonProvider

import (
	jsonIO "encoding/json"
	"errors"
	"golang-shortlink/pkg/shortLinkService/model"
	"io/ioutil"
)

type JsonProvider struct {
	paths model.Paths
	Error error
}

func (p *JsonProvider) LoadPaths(src string) {
	if p.Error != nil {
		return
	}

	var jsonFile model.URLMapping
	file, err := ioutil.ReadFile(src)
	if err != nil {
		p.Error = err
		return
	}

	if err = jsonIO.Unmarshal(file, &jsonFile); err != nil {
		p.Error = err
		return
	}

	var data interface{} = jsonFile.Paths
	var ok bool
	if p.paths, ok = data.(model.Paths); ok != true {
		p.Error = errors.New("json file is invalid")
		return
	}
	p.Error = nil
}

func (p *JsonProvider) GetURL(shortURL string) string {
	if p.Error != nil {
		return ""
	}

	var longURL string
	if longURL = p.paths[shortURL]; longURL == "" {
		p.Error = errors.New("short link not found")
		return ""
	}
	p.Error = nil
	return longURL
}

func (p *JsonProvider) GetErr() error {
	return p.Error
}

func Create() *JsonProvider {
	return &JsonProvider{}
}
