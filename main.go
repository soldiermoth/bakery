package main

import (
	"net/http"
	"os"

	"github.com/cbsinteractive/bakery/config"
	"github.com/cbsinteractive/bakery/handlers"

	"github.com/akrylysov/algnhsa"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		os.Exit(-1)
	}

	logger := c.GetLogger()
	handler := handlers.LoadHandler(c)

	_, isLambda := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")

	// check if it's running on lambda environment or not
	if isLambda {
		algnhsa.ListenAndServe(handler, nil)
	} else {
		logger.Infof("Starting Bakery on %s", c.Listen)
		http.Handle("/", handler)
		http.ListenAndServe(c.Listen, nil)
	}
}
