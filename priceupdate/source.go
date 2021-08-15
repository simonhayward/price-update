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
	Convert  bool   `json:"c,omitempty"`
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

// SetPrice - get response from URL then clean, convert price
func (s *Source) SetPrice() error {
	if err := s.SetReponse(); err != nil {
		return err
	}

	price, err := s.Parse()
	if err != nil {
		return fmt.Errorf("parse failed for: %s error: %s", s.URL, err)
	}

	s.Price = price
	s.CleanPrice()
	s.ConvertPrice()
	return nil
}

// CleanPrice - remove formatting chrs
func (s *Source) CleanPrice() {
	s.Price = strings.Replace(s.Price, ",", "", -1)
}

// ConvertPrice - if set on a price source shift decimal point (pounds to pence)
func (s *Source) ConvertPrice() {
	if s.Convert {
		price := fmt.Sprintf("%s00", s.Price)
		index := strings.Index(price, ".")
		if index != -1 {
			// remove point
			price = fmt.Sprintf("%s%s", price[:index], price[index+1:])
			// shift point
			price = fmt.Sprintf("%s.%s", price[:index+2], price[index+2:])
		}
		s.Price = price
	}
}

// Parse - extract match from response
func (s *Source) Parse() (string, error) {
	r, err := regexp.Compile(s.Pattern)
	if err != nil {
		return "", err
	}

	matches := r.FindStringSubmatch(string(s.Response))
	matchesLen := len(matches)

	if matches == nil || matchesLen == 1 {
		return "", fmt.Errorf("no match for regex: %s in response: %s", s.Pattern, string(s.Response))
	}

	if len(matches) > 2 {
		LogOutput(Log{Message: fmt.Sprintf("more than: %d matches found: %s", len(matches), s.URL), Severity: Warning})
	}

	return matches[1], nil

}
