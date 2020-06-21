package httputil

import (
	"net/http"

	"github.com/mdzio/go-logging"
)

var (
	log = logging.Get("auth-handler")
)

// SingleAuthHandler wraps another http.Handler and forces the specified
// authentication from the HTTP client.
type SingleAuthHandler struct {
	http.Handler

	User     string
	Password string

	// Realm must only contain valid characters for an HTTP header value and no
	// double quotes.
	Realm string
}

func (h *SingleAuthHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	user, passwd, ok := req.BasicAuth()

	// no credentials
	if !ok {
		log.Tracef("Not authenticated: %s", req.RemoteAddr)
		h.sendAuth(rw, req)
		return
	}

	// check credentials
	if user != h.User {
		log.Warningf("Unknown user %s: %s", user, req.RemoteAddr)
		h.sendAuth(rw, req)
		return
	}
	if passwd != h.Password {
		log.Warningf("Invalid password for user %s: %s", user, req.RemoteAddr)
		h.sendAuth(rw, req)
		return
	}

	// credentials ok
	h.Handler.ServeHTTP(rw, req)
}

func (h *SingleAuthHandler) sendAuth(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("WWW-Authenticate", "Basic realm=\""+h.Realm+"\", charset=\"UTF-8\"")
	http.Error(rw, "Unauthorized", http.StatusUnauthorized)
}
