package handlers

import (
	"errors"
	resp "gushort/internal/lib/api/response"
	"gushort/internal/lib/logger/sl"
	"gushort/internal/lib/random"
	"gushort/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type UrlSaver interface {
	Save(url string, reqAlias *string) (string, error)
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			validateErrs := err.(validator.ValidationErrors)
			render.JSON(w, r, resp.ValidationError(validateErrs))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomAlias(aliasLength)
		}

		savedAls, err := urlSaver.Save(req.Url, &alias)
		if errors.Is(err, storage.ErrUrlAlreadyExists) {
			log.Info("url already exists", slog.String("url", req.Url))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.String("alias", savedAls))

		render.JSON(w, r, Response{resp.OK(), alias})
	}
}
