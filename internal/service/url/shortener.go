package url

import (
	"gushort/internal/lib/random"
	"gushort/internal/storage/sqlite"
	"log/slog"
)

type UrlShortener struct {
	log     *slog.Logger
	storage *sqlite.Storage
}

const aliasLength = 6

func New(log *slog.Logger, storage *sqlite.Storage) *UrlShortener {
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
	}

	alias = *reqAlias

	id, err := s.storage.SaveUrl(url, alias)
	if err != nil {
		return "", err
	}

	log.Info("url saved", slog.Int64("id", id))

	return alias, nil
}

func (s *UrlShortener) Get(reqAlias string) (string, error) {
	const op = "service.redirect"

	url, err := s.storage.GetUrlByAlias(reqAlias)
	if err != nil {
		return "", err
	}

	s.log.Info("url saved", slog.String("url", url))

	return url, nil
}
