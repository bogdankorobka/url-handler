package server

import (
	"net/http"

	"github.com/go-chi/render"
)

// @Summary      UrlHandler
// @Tags          UrlHandler
// @Accept       json
// @Produce      json
// @Success      200    "OK"
// @Failure      422    "Validation error"
// @Failure      429    "Too many requests"
// @Failure      500    "Server Error"
// @Router       /url-handler [POST]
func UrlHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://golang.org")
		if err != nil {
			render.Status(r, http.StatusInternalServerError)

			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, render.M{
			"page_size": resp.ContentLength,
		})
	}
}
