package app

import (
	"fmt"
	"golang-shortlink/pkg/shortLinkService/infrastructure/httpServer"
	"golang-shortlink/pkg/shortLinkService/infrastructure/json"
	"golang-shortlink/pkg/shortLinkService/infrastructure/logger"
)

func InitService(src string, port string) error {
	if err := json.LoadPaths(src); err != nil {
		return err
	}

	httpServer.Response = Response{}

	if err := httpServer.StartHttpServer(port); err != nil {
		return err
	}
	return nil
}

type Response struct{}

func (r Response) Handler(path string) string {
	shortUrl, err := json.GetURL(path)
	if err != nil {
		logger.Logger.Error(err, "Error handler")
		return fmt.Sprintf("<h1>Такой ссылки не существует</h1>")
	}
	info := fmt.Sprint("Handler: ", path, " -> ", shortUrl)
	logger.Logger.Info(info)
	return fmt.Sprintf("<script>location='%s';</script>", shortUrl)
	//return shortUrl
}
