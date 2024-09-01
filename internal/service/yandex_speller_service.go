package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/avran02/kode/config"
	"github.com/avran02/kode/internal/dto"
)

type YandexSpellerService interface {
	CheckText(text string) ([]dto.SpellError, error)
}

type yandexSpellerService struct {
	lang         string
	checkTextURL string
	options      string
	client       *http.Client
}

func (s *yandexSpellerService) CheckText(text string) ([]dto.SpellError, error) {
	params := url.Values{}
	params.Add("text", text)
	params.Add("lang", s.lang)
	params.Add("options", s.options)

	resp, err := s.client.Get(s.checkTextURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	var errors []dto.SpellError
	if err := json.NewDecoder(resp.Body).Decode(&errors); err != nil {
		return nil, err
	}

	return errors, nil
}

func NewYandexSpellerService(config config.YandexSpeller) YandexSpellerService {
	checkTextURL, err := url.JoinPath(config.URL, "checkText")
	if err != nil {
		log.Fatal("failed to create YandexSpellerService", "error", err.Error())
	}
	return &yandexSpellerService{
		lang:         config.Language,
		checkTextURL: checkTextURL,
		options:      config.Options,
		client:       &http.Client{},
	}
}
