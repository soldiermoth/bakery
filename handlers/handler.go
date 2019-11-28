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

		masterManifestPath, mediaFilters, err := parsers.URLParse(r.URL.Path)
		if err != nil {
			httpError(c, w, err, "failed parsing url", http.StatusInternalServerError)
			return
		}

		client := c.Client.New()
		manifestURL := c.OriginHost + masterManifestPath
		resp, err := client.Get(manifestURL)
		if err != nil {
			httpError(c, w, err, "failed fetching origin url", http.StatusInternalServerError)
			return
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
			httpError(c, w, err, "failed to filter manifest", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, filteredManifest)
	})
}

func httpError(c config.Config, w http.ResponseWriter, err error, message string, code int) {
	logger := c.GetLogger()
	logger.WithError(err).Infof(message)
	http.Error(w, message+": "+err.Error(), code)
}
