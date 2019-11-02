package handlers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cbsinteractive/bakery/config"
	"github.com/cbsinteractive/bakery/filters"
	"github.com/cbsinteractive/bakery/parsers"
)

// LoadHandler loads the handler for all the requests
func LoadHandler(c config.Config) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger := c.GetLogger()

		mediaFilters, err := parsers.URLParse(r.URL.Path)
		if err != nil {
			logger.WithError(err).Fatal(w, "failed parsing url")
		}

		logger.Infof("Parsed url with ", mediaFilters)

		client := c.Client.New()
		manifestURL := c.OriginHost + r.URL.Path
		resp, err := client.Get(manifestURL)
		if err != nil {
			logger.WithError(err).Fatal(w, "failed fetching url")
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		var f filters.Filter
		if mediaFilters.Protocol == parsers.ProtocolHLS {
			f = filters.NewHLSFilter(buf.String(), c)
		} else if mediaFilters.Protocol == parsers.ProtocolDASH {
			f = filters.NewDASHFilter(buf.String(), c)
		}

		filteredManifest, err := f.FilterManifest(mediaFilters)
		if err != nil {
			logger.WithError(err).Fatal(w, "failed to filter")
		}

		fmt.Fprintf(w, filteredManifest)
	})
}
