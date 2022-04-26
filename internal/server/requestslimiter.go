package server

import (
	"net/http"

	"github.com/go-chi/render"
)

type requestsLimiter struct {
	limit uint
	sema  chan struct{}
}

func Limit(requestsLimit uint) func(next http.Handler) http.Handler {
	rl := &requestsLimiter{
		limit: requestsLimit + 1,
		sema:  make(chan struct{}, requestsLimit+1),
	}

	return rl.Handler
}

func (l *requestsLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.sema <- struct{}{}
		defer func() {
			<-l.sema
		}()

		if len(l.sema) == int(l.limit) {
			render.Status(r, http.StatusTooManyRequests)
			render.JSON(w, r, render.M{
				"message": "try again later",
			})

			return
		}

		next.ServeHTTP(w, r)
	})
}
