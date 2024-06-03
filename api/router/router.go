// Defines API routes and serves route handlers using the Router struct.
package router

import (
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/agnate/qlikrestapi/api/entity/message"
	myCtx "github.com/agnate/qlikrestapi/internal/context"
)

// Contains the routers for API to serve.
type Router struct {
	routes []route
}

// Contains a route for use in Router.
type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

// Build a new Router containing all of the API routes and handlers.
func New(db *sql.DB) *Router {
	msgAPI := message.New(db)

	return &Router{
		routes: []route{
			newRoute(http.MethodGet, "/api/v1/messages", msgAPI.List),
		},
	}
}

// Build a new route to store in Router. Supports REGEX in the pattern.
//
// # Parameters
//   - method: Use constants from http package (ex: [net/http.MethodGet], [net/http.MethodPost])
//   - pattern: Supports REGEX, with pattern being sandwiched between start/end metachars as follows: ^pattern$
//   - handler: Function to invoke when router is matched
func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

// Create new http.Handler for this Router for use by [net/http.ListenAndServe].
func (rt *Router) NewHandler() http.Handler {
	return http.HandlerFunc(rt.serve)
}

// Uses the [net/http.Request] to match a valid, allowed route and invoke its handler.
func (rt *Router) serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range rt.routes {
		// Use regex to match the route and store the match in the Context.
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			// Check that request and route match for their GET/POST method.
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			// Add the regex match to the Context and invoke the route's handler.
			ctx := myCtx.SetContextRoute(r.Context(), matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	// We found a matching route, but the methods didn't match (GET/POST), so
	// inform user that method isn't allowed.
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, newHttpStatusMsg(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

// statusCode: Use constants from http package (ex: [net/http.StatusMethodNotAllowed])
func newHttpStatusMsg(statusCode int) string {
	return strconv.Itoa(statusCode) + " " + http.StatusText(statusCode)
}