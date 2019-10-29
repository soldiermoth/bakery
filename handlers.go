package main

import (
	"fmt"
	"net/http"
)

// LoadHandler loads the handler for all the requests
func LoadHandler(c Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		fmt.Fprintf(w, "It works!")
	})
}
