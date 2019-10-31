package main

import (
	"bytes"
	"fmt"
	"net/http"
)

// LoadHandler loads the handler for all the requests
func LoadHandler(c Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		manifestURL := "https://vod-gcs-cedexis.cbsaavideo.com/intl_vms/2019/10/17/1625066563842/123002_cenc_dash/stream.mpd"

		client := c.Client.New()
		resp, err := client.Get(manifestURL)
		if err != nil {
			fmt.Fprintf(w, "failed fetching url")
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Fprintf(w, buf.String())
	})
}
