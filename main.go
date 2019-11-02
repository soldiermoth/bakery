package main

import (
	"net"
	"net/http"
	"os"

	"github.com/cbsinteractive/bakery/config"
	"github.com/cbsinteractive/bakery/handlers"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		os.Exit(-1)
	}

	logger := c.GetLogger()

	listener, err := net.Listen("tcp", c.Listen)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize listener")
	}

	logger.Infof("Starting Bakery on %s", listener.Addr())
	handler := handlers.LoadHandler(c)
	err = http.Serve(listener, handler)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize server")
	}
}
