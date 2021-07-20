package json

import (
	jsonIO "encoding/json"
	"errors"
	"golang-shortlink/pkg/shortLinkService/model"
	"io/ioutil"
)

var paths interface{}

func LoadPaths(src string) error {
	var jsonFile model.URLMapping
	file, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if err = jsonIO.Unmarshal(file, &jsonFile); err != nil {
		return err
	}

	paths = jsonFile.Paths
	if _, ok := paths.(model.Paths); ok != true {
		return errors.New("Paths not found")
	}
	return nil
}

func GetURL(longUrl string) (string, error) {
	var shortUrl string
	if val, ok := paths.(model.Paths); ok != true {
		return "", errors.New("Paths not found")
	} else {
		if shortUrl = val[longUrl]; shortUrl == "" {
			return "", errors.New("Short link not found")
		}
	}
	return shortUrl, nil
}
