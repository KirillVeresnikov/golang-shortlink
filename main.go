package main

import (
	"flag"
	"golang-shortlink/pkg/shortLinkService/app"
	"golang-shortlink/pkg/shortLinkService/infrastructure/logger"
	"os"
)

func main() {
	src := "./paths.json"
	file := flag.Bool("f", false, "a bool")
	flag.Parse()
	if file != nil && *file {
		src = os.Args[1]
	}

	err := app.InitService(src, ":8000")
	if err != nil {
		logger.Logger.Error(err, "Application stop")
	}
}
