package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

var sensorManager *SensorManager
var detectionManager *DetectionManager

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

	centroidLatitudeValue := os.Getenv("CENTROID_LATITUDE")
	if centroidLatitudeValue == "" {
		log.Fatal("environment variable CENTROID_LATITUDE cannot be empty")
	}
	centroidLatitude, err := strconv.ParseFloat(centroidLatitudeValue, 64)
	if err != nil {
		log.Fatal("cannot convert environment variable CENTROID_LATITUDE to float:", err)
	}

	centroidLongitudeValue := os.Getenv("CENTROID_LONGITUDE")
	if centroidLongitudeValue == "" {
		log.Fatal("environment variable CENTROID_LONGITUDE cannot be empty")
	}
	centroidLongitude, err := strconv.ParseFloat(centroidLongitudeValue, 64)
	if err != nil {
		log.Fatal("cannot convert environment variable CENTROID_LONGITUDE to float:", err)
	}

	centroid := Point{Latitude: centroidLatitude, Longitude: centroidLongitude}

	api := NewWildlifeNLAPI(apiURL, token)
	sensorManager = NewSensorManager(centroid, numberOfSensors, api)
	species, err := api.GetSpecies()
	if err != nil {
		log.Fatal("cannot get species from api:", err)
	}
	detectionManager = NewDetectionManager(centroid, species, api)

	update()
	for range time.Tick(time.Duration(interval) * time.Minute) {
		update()
	}
}

func update() {
	if err := sensorManager.Update(); err != nil {
		log.Println("ERROR: could not update sensors:", err)
	}
	if err := detectionManager.Update(); err != nil {
		log.Println("ERROR: could not update detections:", err)
	}
}
