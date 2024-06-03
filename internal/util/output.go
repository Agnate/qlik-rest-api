package util

import (
	"log"
	"net/http"
	"strconv"
)

// statusCode: Use constants from http package (ex: [net/http.StatusMethodNotAllowed])
func NewHttpStatusMsg(statusCode int) string {
	return strconv.Itoa(statusCode) + " " + http.StatusText(statusCode)
}

// Used when validation or data loading fails for an endpoint and we want a consistent
// output displayed to our users. Errors will be logged.
func NoAPIEndpoint(w http.ResponseWriter, err error) {
	// TODO: Add logging for invalid endpoints in case we need to monitor spammers.
	log.Println(err)
	http.Error(w, NewHttpStatusMsg(http.StatusNotFound), http.StatusNotFound)
}
