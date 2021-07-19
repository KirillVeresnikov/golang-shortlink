package logger

import (
	"fmt"
	"time"
)

var Logger log

type log struct{}

func (l log) Info(data interface{}) {
	fmt.Println(time.Now(), ": INFO: ", data)
}

func (l log) Error(err error, data interface{}) {
	fmt.Println(time.Now(), ": ERROR: ", err, "; ", data)
}
