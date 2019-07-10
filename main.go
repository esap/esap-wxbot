package main

import (
	"flag"
	"fmt"
	// "log"
	// "net/http"
)

type ServiceError struct {
	Msg string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

func Error(msg string) error {
	return &ServiceError{msg}
}

const (
	Version = "4.0.1"
)

var debug = flag.String("d", "off", "if on debug mode")

func main() {
	flag.Parse()
	fmt.Println(Version)
	fmt.Printf("debug mode %s\n", *debug)

	wx := wxweb{}
	wx.start()
}
