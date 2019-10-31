package main

import (
	"net"
	"net/http"
	"os"
)

func main() {
	c, err := LoadConfig()
	if err != nil {
		os.Exit(-1)
	}

	logger := c.GetLogger()

	listener, err := net.Listen("tcp", c.Listen)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize listener")
	}

	logger.Infof("Starting Bakery on %s", listener.Addr())
	handler := LoadHandler(c)
	err = http.Serve(listener, handler)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize server")
	}
}
