package handlers

import (
	resp "gushort/internal/lib/api/response"
	"gushort/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Url string `json:"url,omitempty"`
}

type UrlGetter interface {
	Get(reqAlias string) (string, error)
}

func New(log *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("empty alias url param error")

			render.JSON(w, r, resp.Error("empty alias url param error"))
			return
		}

		url, err := urlGetter.Get(alias)
		if err != nil {
			log.Error("failed to get url by alias", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get url by alias"))
			return
		}

		log.Info("url got", slog.String("alias", alias), slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
