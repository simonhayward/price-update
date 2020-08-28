package priceupdate

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Source struct {
	URL      string `json:"url"`
	Pattern  string `json:"p"`
	Response []byte
	Price    string
}

type Sources []Source

type Security struct {
	ISIN    string  `json:"isin"`
	Sources Sources `json:"src"`
}

type Securities []Security

// GetSecurities - from url convert into securities
func GetSecurities(url string) (*Securities, error) {
	body, err := GetResponse(url)
	if err != nil {
		return nil, err
	}
	return GenerateSecurities(body)
}

// GenerateSecurities - from json blob convert into Securities
func GenerateSecurities(body []byte) (*Securities, error) {
	s := Securities{}
	err := json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// SetReponse - response from source URL
func (s *Source) SetReponse() error {
	var err error
	s.Response, err = GetResponse(s.URL)
	if err != nil {
		return err
	}

	return nil
}

func (s *Source) SetPrice() error {
	if err := s.SetReponse(); err != nil {
		return err
	}

	price, err := s.Parse()
	if err != nil {
		return fmt.Errorf("parse failed for: %s error: %s", s.URL, err)
	}

	// Cleanup price
	s.Price = strings.Replace(price, ",", "", -1)

	return nil
}

func (s *Source) Parse() (string, error) {
	r, err := regexp.Compile(s.Pattern)
	if err != nil {
		return "", err
	}

	matches := r.FindStringSubmatch(string(s.Response))
	matchesLen := len(matches)

	if matches == nil || matchesLen == 1 {
		return "", fmt.Errorf("no match for regex: %s", s.Pattern)
	}

	if len(matches) > 2 {
		fmt.Println(fmt.Sprintf(`{"message": "more than: %d matches found: %s", "severity": "warning"}`, len(matches), s.URL))
	}

	return matches[1], nil

}
