package httpServer

import (
	"fmt"
	"golang-shortlink/pkg/shortLinkService/infrastructure/logger"
	"net/http"
)

type HttpHandler struct{}
type Handler interface {
	Handler(path string) string
}

var Response Handler

func StartHttpServer(port string) error {
	handle := HttpHandler{}

	if err := http.ListenAndServe(port, handle); err != nil {
		logger.Logger.Error(err, "")
		return err
	}
	logger.Logger.Info("Service is starting")
	return nil
}

func (h HttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, Response.Handler(request.URL.Path))
}
