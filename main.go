package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

const addr = ":8080"

var centralPoint = Point{Latitude: 52.088644, Longitude: 5.172076} // Botanische Tuinen

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		log.Fatal("environment variable API_URL cannot be empty")
	}

	token := os.Getenv("BEARER_TOKEN")
	if token == "" {
		log.Fatal("environment variable BEARER_TOKEN cannot be empty")
	}

	intervalValue := os.Getenv("INTERVAL")
	if intervalValue == "" {
		log.Fatal("environment variable INTERVAL cannot be empty")
	}
	interval, err := strconv.Atoi(intervalValue)
	if err != nil {
		log.Fatal("cannot convert environment variable INTERVAL to integer:", err)
	}

	numberOfSensorsValue := os.Getenv("NUM_SENSORS")
	if numberOfSensorsValue == "" {
		log.Fatal("environment variable NUM_SENSORS cannot be empty")
	}
	numberOfSensors, err := strconv.Atoi(numberOfSensorsValue)
	if err != nil {
		log.Fatal("cannot convert environment variable NUM_SENSORS to integer:", err)
	}
	if numberOfSensors < 1 || numberOfSensors > 99 {
		log.Fatal("Environment variable NUM_SENSORS must be 0 < x < 100.")
	}

	api := NewWildlifeNLAPI(apiURL, token)
	go func() {
		webServer := NewWebServer(addr, api)
		log.Println(webServer.ListenAndServe())
	}()

	sensorManager := NewSensorManager(numberOfSensors, api)
	tick := time.Tick(time.Duration(interval) * time.Minute)
	for range tick {
		if err := sensorManager.Update(); err != nil {
			log.Print("ERROR: could not update sensors:", err)
		}
	}
}
