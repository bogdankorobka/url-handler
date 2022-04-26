package server

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
)

type UrlListRequest struct {
	UrlList []string `json:"url_list" validate:"required,gt=0,lte=100,dive,required,url" example:"https://google.com,https://yandex.ru"`
}

// @Summary      UrlHandler
// @Tags         UrlHandler
// @Accept       json
// @Produce      json
// @Param Request body UrlListRequest true "Url List"
// @Success      200    "OK"
// @Failure      422    "Validation error"
// @Failure      429    "Too many requests"
// @Failure      500    "Server Error"
// @Router       /url-handler [POST]
func UrlHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &UrlListRequest{}

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

		message := make(map[string]int64, len(req.UrlList))
		for _, url := range req.UrlList {
			httpResp, err := http.Get(url)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)

				return
			}

			cl := httpResp.ContentLength

			if cl == -1 {
				content, err := ioutil.ReadAll(httpResp.Body)
				if err != nil {
					render.Status(r, http.StatusInternalServerError)

					return
				}
				defer httpResp.Body.Close()

				cl = int64(len(content))
			}

			message[url] = cl
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, message)
	}
}
