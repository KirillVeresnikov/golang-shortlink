package jsonProvider

import (
	jsonIO "encoding/json"
	"errors"
	"golang-shortlink/pkg/shortLinkService/model"
	"io/ioutil"
)

type JsonProvider struct {
	paths model.Paths
	Err   error
}

func (p *JsonProvider) LoadPaths(src string) {
	if p.Err != nil {
		return
	}

	var jsonFile model.URLMapping
	file, err := ioutil.ReadFile(src)
	if err != nil {
		p.Err = err
		return
	}

	if err = jsonIO.Unmarshal(file, &jsonFile); err != nil {
		p.Err = err
		return
	}

	var data interface{} = jsonFile.Paths
	var ok bool
	if p.paths, ok = data.(model.Paths); ok != true {
		p.Err = errors.New("json file is invalid")
		return
	}
	p.Err = nil
}

func (p *JsonProvider) GetURL(shortURL string) string {
	if p.Err != nil {
		return ""
	}

	var longURL string
	if longURL = p.paths[shortURL]; longURL == "" {
		p.Err = errors.New("short link not found")
		return ""
	}
	p.Err = nil
	return longURL
}

func (p *JsonProvider) GetErr() error {
	defer func(p *JsonProvider) { p.Err = nil }(p)
	return p.Err
}

func Create() *JsonProvider {
	return &JsonProvider{}
}
