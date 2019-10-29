package main

import (
	"net"
	"net/http"
)

func main() {
	c := LoadConfig()
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
