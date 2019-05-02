package fairway

import (
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logger.Infof(
			"%s\t%s\t%s",
			time.Now().UTC().Format("2006-01-02 15:04:05.999Z"),
			r.Method,
			r.RequestURI,
		)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
