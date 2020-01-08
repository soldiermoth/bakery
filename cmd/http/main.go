package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/cbsinteractive/bakery/pkg/config"
	"github.com/cbsinteractive/bakery/pkg/handlers"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
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
		if err := http.ListenAndServe(c.Listen, nil); err != nil {
			log.Fatal(err)
		}
	}
}
