package app

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type Service struct {
	rnd    *rand.Rand
	urlDAO *UrlDAO
}

func NewService(urlDAO *UrlDAO) *Service {
	return &Service{
		rnd:    rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0)),
		urlDAO: urlDAO,
	}
}

func (s *Service) Shorten(ctx context.Context, url string, ttlDays int) (*ShortURL, error) {
	shortURL := &ShortURL{
		URL:      url,
		ExpireAt: getExpirationTime(ttlDays),
	}
	for range 10 {
		shortURL.ID = s.generateRandomID()
		err := s.urlDAO.Insert(ctx, shortURL)
		if err == nil {
			return shortURL, nil
		}
		if !mongo.IsDuplicateKeyError(err) {
			return nil, err
		}
	}
	return nil, errors.Errorf("failed to create short link due to collision")
}

func (s *Service) Update(ctx context.Context, id string, url string, ttlDays int) (*ShortURL, error) {
	shortURL, err := s.urlDAO.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	shortURL.URL = url
	shortURL.ExpireAt = getExpirationTime(ttlDays)
	return shortURL, s.urlDAO.Update(ctx, shortURL)
}

func (s *Service) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	sURL, err := s.urlDAO.FindByID(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return sURL.URL, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.urlDAO.DeleteByID(ctx, id)
}

func (s *Service) generateRandomID() string {
	const idLenght = 6
	id := make([]rune, idLenght)
	for i := range id {
		id[i] = symbols[s.rnd.IntN(len(symbols))]
	}
	return string(id)
}

func getExpirationTime(ttlDays int) *time.Time {
	if ttlDays <= 0 {
		return nil
	}
	t := time.Now().Add(24 * time.Hour * time.Duration(ttlDays))
	return &t
}
