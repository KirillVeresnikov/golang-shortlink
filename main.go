package main

import (
	"context"
	"flag"
	"fmt"
	"golang-shortlink/pkg/shortLinkService/app"
	"net/http"
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
		s.server.ListenAndServe()
	}()
}

func (s *httpServer) shutdownServer() {
	s.server.Shutdown(context.Background())
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
		fmt.Fprintln(w, "Введите короткую ссылку")
	} else {
		longURL := service.GetLongURL(URL)
		if err := service.GetErr(); err == nil {
			//http.Redirect(w, r, longURL, 302)
			fmt.Fprintln(w, longURL)
		} else {
			fmt.Fprintln(w, "Error: ", err)
		}
	}
}

var service Service
var server httpServer

func main() {
	src := "./paths.json"
	file := flag.Bool("f", false, "a file path")
	flag.Parse()
	if file != nil && *file {
		src = os.Args[1]
	}

	service = app.Create(src)
	server = httpServer{}
	server.handler = HttpHandler{}

	server.startServer()
	server.getKillSignalChan()
	server.waitKillSignalChan()
	server.shutdownServer()
}
