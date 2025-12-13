package middlewares

import (
	"net/http"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

func ShouldJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "application/json" && r.Method != "GET" && r.Method != "PATCH" && r.Method != "DELETE" {
			utils.WriteJSONError(w, http.StatusUnsupportedMediaType, "unsupported media type")
			return
		}

		next.ServeHTTP(w, r)
	})
}
