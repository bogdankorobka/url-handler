package server

import (
	"net/http"

	"github.com/bogdankorobka/url-handler/internal/services"
	"github.com/go-chi/render"
)

type URLListRequest struct {
	URLList []string `json:"url_list" validate:"required,gt=0,lte=100,dive,required,url" example:"https://google.com,https://yandex.ru"`
}

// @Summary      UrlHandler
// @Tags         UrlHandler
// @Accept       json
// @Produce      json
// @Param Request body URLListRequest true "Url List"
// @Success      200    "OK"
// @Failure      422    "Validation error"
// @Failure      429    "Too many requests"
// @Failure      500    "Server Error"
// @Router       /url-handler [POST]
func UrlHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const numWorkers = 3

		req := &URLListRequest{}

		// decode body to request
		err := render.DecodeJSON(r.Body, req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, render.M{"message": err.Error()})

			return
		}

		// validate request
		valErrs, err := Validate(req)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, render.M{"message": "server error"})

			return
		}

		if valErrs != nil {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, valErrs)

			return
		}

		URLService := services.NewURLServie(numWorkers, req.URLList)

		// work service
		res, ok := URLService.Start(r.Context())
		if !ok {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, render.M{"message": "url unavailable"})

			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, res)
	}
}
