package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type WildlifeNLAPI struct {
	url   string
	token string
}

func NewWildlifeNLAPI(url string, token string) *WildlifeNLAPI {
	a := WildlifeNLAPI{url: url, token: token}
	return &a
}

func (a *WildlifeNLAPI) SendReading(sensor Sensor) error {
	data, err := json.Marshal(sensor)
	if err != nil {
		return fmt.Errorf("could not json marshal: %w", err)
	}
	r, _ := http.NewRequest(http.MethodPost, a.url+"/borne-sensor-reading/", bytes.NewReader(data))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/problem+json")
	r.Header.Add("Authorization", "Bearer "+a.token)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("failed sending reading: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 204 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed sending reading, got status code: "+strconv.Itoa(res.StatusCode), string(body))
	}
	return nil
}

func (a *WildlifeNLAPI) GetAnimals() ([]byte, error) {
	r, _ := http.NewRequest(http.MethodGet, a.url+"/animals/", nil)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/problem+json")
	r.Header.Add("Authorization", "Bearer "+a.token)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed getting animals: %w", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed getting animals, got status code: "+strconv.Itoa(res.StatusCode), string(body))
	}
	return body, nil
}
