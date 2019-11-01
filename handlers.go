package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cbsinteractive/go-dash/mpd"
)

// LoadHandler loads the handler for all the requests
func LoadHandler(c Config) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger := c.GetLogger()

		manifestURL := c.OriginHost + r.URL.Path
		//mediaFilters, err := Parser(manifestURL)
		//if err != nil {
		//logger.WithError(err).Fatal(w, "failed parsing url")
		//}

		client := c.Client.New()
		resp, err := client.Get(manifestURL)
		if err != nil {
			logger.WithError(err).Fatal(w, "failed fetching url")
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		manifest, err := mpd.ReadFromString(buf.String())
		if err != nil {
			logger.WithError(err).Fatal(w, "failed to parse mpd")
		}

		newManifest, err := manifest.WriteToString()
		if err != nil {
			logger.WithError(err).Fatal(w, "failed to generate mpd")
		}

		fmt.Fprintf(w, newManifest)
	})
}
