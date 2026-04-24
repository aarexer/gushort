package url

import (
	"fmt"
	"gushort/internal/lib/random"
	"log/slog"
)

type Storage interface {
	SaveUrl(url string, alias string) (int64, error)
	GetUrlByAlias(alias string) (string, error)
}

type UrlShortener struct {
	log     *slog.Logger
	storage Storage
}

const aliasLength = 6

func New(log *slog.Logger, storage Storage) *UrlShortener {
	return &UrlShortener{log: log, storage: storage}
}

func (s *UrlShortener) Save(url string, reqAlias *string) (string, error) {
	const op = "service.save"

	log := s.log.With(
		slog.String("op", op),
	)

	var alias string
	if reqAlias == nil || *reqAlias == "" {
		alias = random.NewRandomAlias(aliasLength)
	} else {
		alias = *reqAlias
	}

	id, err := s.storage.SaveUrl(url, alias)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("url saved", slog.Int64("id", id))

	return alias, nil
}

func (s *UrlShortener) Get(reqAlias string) (string, error) {
	const op = "service.redirect"

	url, err := s.storage.GetUrlByAlias(reqAlias)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("url retrieved", slog.String("url", url))

	return url, nil
}
