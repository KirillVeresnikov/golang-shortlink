package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang-shortlink/pkg/shortLinkService/app"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

type httpServer struct {
	server         *http.Server
	handler        http.Handler
	killSignalChan chan os.Signal
	err            error
}

func (s *httpServer) startServer() {
	s.server = &http.Server{Addr: ":8000", Handler: s.handler}
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err.Error() != errors.New("http: Server closed").Error() {
			logger.LogFatal("starting fail", err)
		}
	}()
}

func (s *httpServer) shutdownServer() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		logger.LogFatal("shutdown failed", err)
	}
}

func (s *httpServer) getKillSignalChan() {
	s.killSignalChan = make(chan os.Signal, 1)
	signal.Notify(s.killSignalChan, os.Interrupt, syscall.SIGTERM)
}

func (s *httpServer) waitKillSignalChan() {
	func(killSignalChan <-chan os.Signal) {
		killSignal := <-killSignalChan
		switch killSignal {
		case os.Interrupt:
			fmt.Println("Interrupt")
		case syscall.SIGTERM:
			fmt.Println("SIGTERM")
		}
	}(s.killSignalChan)
}

type Service interface {
	GetLongURL(shortURL string) string
	GetErr() error
}

type HttpHandler struct{}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if URL := r.URL.Path; URL == "/" {
		_, wErr := fmt.Fprintln(w, "Введите короткую ссылку")
		if wErr != nil {
			logger.LogFatal("response writer error", wErr)
		}
	} else {
		longURL := service.GetLongURL(URL)
		if err := service.GetErr(); err == nil {
			http.Redirect(w, r, longURL, http.StatusFound)
			logger.LogServerError(URL, "successful")
		} else {
			logger.LogError(r.URL, err)
			_, wErr := fmt.Fprintln(w, "Короткая ссылка не найдена")
			if wErr != nil {
				logger.LogFatal("response writer error", wErr)
			}
		}
	}
}

var service Service
var server httpServer
var logger serverErrorLogger

func main() {
	log := initializeLogger("debug")
	logger = serverErrorLogger{log}

	src := "./paths.json"
	file := flag.Bool("f", false, "a file path")
	flag.Parse()
	if file != nil && *file {
		if len(os.Args) >= 3 {
			src = os.Args[2]
		} else {
			src = ""
		}
	}

	service = app.Create(src)
	if err := service.GetErr(); err != nil {
		logger.LogFatal("create service failed", err)
	}

	server = httpServer{}
	server.handler = HttpHandler{}

	server.startServer()
	logger.LogInfo("starting complete")
	server.getKillSignalChan()
	server.waitKillSignalChan()
	server.shutdownServer()
}

func initializeLogger(logLevelParam string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	logLevel, err := logrus.ParseLevel(logLevelParam)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	logger.SetLevel(logLevel)

	return logger
}

type serverErrorLogger struct {
	Logger *logrus.Logger
}

func (logger *serverErrorLogger) LogServerError(address string, info string) {
	logger.Logger.WithFields(logrus.Fields{
		"address": address,
	}).Info(info)
}

func (logger *serverErrorLogger) LogError(requestURL *url.URL, err error) {
	logger.Logger.WithFields(logrus.Fields{
		"err": err.Error(),
		"url": requestURL.String(),
	}).Error("http request failed")
}

func (logger *serverErrorLogger) LogInfo(info string) {
	logger.Logger.Info(info)
}

func (logger *serverErrorLogger) LogFatal(info string, err error) {
	logger.Logger.WithFields(logrus.Fields{
		"err": err.Error(),
	}).Fatal(info)
}
