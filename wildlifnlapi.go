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

type Species struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

func (a *WildlifeNLAPI) GetSpecies() ([]Species, error) {
	data, err := a.getAll("species")
	if err != nil {
		return nil, err
	}
	species := make([]Species, 0)
	if err := json.Unmarshal(data, &species); err != nil {
		return nil, err
	}
	return species, nil
}

func (a *WildlifeNLAPI) getAll(endpoint string) ([]byte, error) {
	r, _ := http.NewRequest(http.MethodGet, a.url+"/"+endpoint+"/", nil)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/problem+json")
	r.Header.Add("Authorization", "Bearer "+a.token)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed getting "+endpoint+": %w", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed getting "+endpoint+", got status code: "+strconv.Itoa(res.StatusCode), string(body))
	}
	return body, nil
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

func (a *WildlifeNLAPI) SendDetection(detector Detector) error {
	data, err := json.Marshal(detector)
	if err != nil {
		return fmt.Errorf("could not json marshal: %w", err)
	}
	r, _ := http.NewRequest(http.MethodPost, a.url+"/detection/", bytes.NewReader(data))
	r.Header.Add("Content-Type", "application/json")
	//r.Header.Add("Accept", "application/problem+json")
	r.Header.Add("Authorization", "Bearer "+a.token)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("failed sending detection: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed sending detection, got status code: "+strconv.Itoa(res.StatusCode), string(body))
	}
	return nil
}
