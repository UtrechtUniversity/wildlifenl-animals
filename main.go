package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var sensors []Sensor
var centralPoint = Point{Latitude: 51.3073247102349, Longitude: 5.658016097401846}

type Sensor struct {
	ID        string    `json:"sensorID"`
	Timestamp time.Time `json:"timestamp"`
	Location  Point     `json:"location"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var (
	apiURL string
	token  string
)

func main() {
	apiURL = os.Getenv("API_URL")
	if apiURL == "" {
		log.Fatal("environment variable API_URL cannot be empty")
	}

	token = os.Getenv("BEARER_TOKEN")
	if token == "" {
		log.Fatal("environment variable BEARER_TOKEN cannot be empty")
	}

	interval_str := os.Getenv("INTERVAL")
	if interval_str == "" {
		log.Fatal("environment variable INTERVAL cannot be empty")
	}
	interval, err := strconv.Atoi(interval_str)
	if err != nil {
		log.Fatal("cannot convert environment variable INTERVAL to integer:", err)
	}

	createSensors()

	tick := time.Tick(time.Duration(interval) * time.Minute)
	for range tick {
		if err := updateSensors(); err != nil {
			log.Print("ERROR: could not update sensors:", err)
		}
	}
}

func createSensors() {
	sensors = make([]Sensor, 0)
	for i := 1; i <= 25; i++ {
		number := strconv.Itoa(i)
		if len(number) == 1 {
			number = "0" + number
		}
		sensors = append(sensors, Sensor{ID: "Sim-Sensor-" + number, Timestamp: time.Now(), Location: centralPoint})
	}
}

func updateSensors() error {
	for _, s := range sensors {
		s.Timestamp = time.Now()
		latDelta := (float64(rand.Intn(7)-3) * .001)
		lonDelta := (float64(rand.Intn(7)-3) * .001)
		s.Location.Latitude = math.Round((s.Location.Latitude+latDelta)*100000) / 100000
		s.Location.Longitude = math.Round((s.Location.Longitude+lonDelta)*100000) / 100000
		if err := sendReading(s); err != nil {
			return err
		}
	}
	return nil
}

func sendReading(sensor Sensor) error {
	data, err := json.Marshal(sensor)
	if err != nil {
		return fmt.Errorf("could not json marshal: %w", err)
	}
	r, _ := http.NewRequest(http.MethodPost, apiURL+"/borne-sensor-reading/", bytes.NewReader(data))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/problem+json")
	r.Header.Add("Authorization", "Bearer "+token)
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
