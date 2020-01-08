package handlers

import (
	"bytes"
	"fmt"
	"github.com/cbsinteractive/bakery/pkg/parsers"
	"net/http"

	"github.com/cbsinteractive/bakery/pkg/config"
	"github.com/cbsinteractive/bakery/pkg/filters"
)

// LoadHandler loads the handler for all the requests
func LoadHandler(c config.Config) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/favicon.ico" {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		defer r.Body.Close()
		logger := c.GetLogger()
		logger.Infof("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)

		// parse all the filters from the URL
		masterManifestPath, mediaFilters, err := parsers.URLParse(r.URL.Path)
		if err != nil {
			httpError(c, w, err, "failed parsing url", http.StatusInternalServerError)
			return
		}

		// request the origin URL
		manifestURL := c.OriginHost + masterManifestPath
		manifestContent, err := fetchManifest(c, manifestURL)
		if err != nil {
			httpError(c, w, err, "failed fetching origin url", http.StatusInternalServerError)
			return
		}

		// create filter associated to the protocol and set
		// response headers accordingly
		var f filters.Filter
		if mediaFilters.Protocol == parsers.ProtocolHLS {
			f = filters.NewHLSFilter(manifestURL, manifestContent, c)
			w.Header().Set("Content-Type", "application/x-mpegURL")
		} else if mediaFilters.Protocol == parsers.ProtocolDASH {
			f = filters.NewDASHFilter(manifestURL, manifestContent, c)
			w.Header().Set("Content-Type", "application/dash+xml")
		} else {
			err := fmt.Errorf("unsupported protocol %q", mediaFilters.Protocol)
			httpError(c, w, err, "failed to select filter", http.StatusBadRequest)
			return
		}

		// apply the filters to the origin manifest
		filteredManifest, err := f.FilterManifest(mediaFilters)
		if err != nil {
			httpError(c, w, err, "failed to filter manifest", http.StatusInternalServerError)
			return
		}

		// write the filtered manifest to the response
		fmt.Fprintf(w, filteredManifest)
	})
}

func fetchManifest(c config.Config, manifestURL string) (string, error) {
	client := c.Client.New()
	resp, err := client.Get(manifestURL)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String(), nil
}

func httpError(c config.Config, w http.ResponseWriter, err error, message string, code int) {
	logger := c.GetLogger()
	logger.WithError(err).Infof(message)
	http.Error(w, message+": "+err.Error(), code)
}
