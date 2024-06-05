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

// 404 Not Found - Used when validation or data loading fails for an endpoint
// and we want a consistent output displayed to our users. Errors will be logged.
func Status404NoAPIEndpoint(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: Add logging for invalid endpoints in case we need to monitor spammers.
	log.Println(err)
	http.NotFound(w, r)
}

// 405 Not Allowed - Used when an unavailable request METHOD is supplied for a route.
// Errors will be logged.
func Status405APINotAllowed(w http.ResponseWriter, err error) {
	// TODO: Add logging for invalid endpoints in case we need to monitor spammers.
	log.Println(err)
	http.Error(w, NewHttpStatusMsg(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// 500 Internal Server Error - Used when an unexpected error occurs and we want a consistent
// output displayed to our users. Errors will be logged.
func Status500APIError(w http.ResponseWriter, err error) {
	// TODO: Add logging for invalid endpoints in case we need to monitor spammers.
	log.Println(err)
	http.Error(w, NewHttpStatusMsg(http.StatusInternalServerError), http.StatusInternalServerError)
}

// 200 Ok - Used to output successful API response.
func Status200APIOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// 201 CREATE - Used to output successful API response.
func Status201APICreate(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

// Writes out content-type header for JSON.
func APIJsonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
