package users

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func AuthenticationMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		if IsAuthenticated(r.Header) {
			h(w, r, params)
		} else {
			http.Error(w, "You are not authenticated", http.StatusUnauthorized)
			return
		}
	}
}
