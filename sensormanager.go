package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type SensorManager struct {
	api     *WildlifeNLAPI
	sensors []Sensor
}

func NewSensorManager(numberOfSensors int, api *WildlifeNLAPI) *SensorManager {
	sensors := make([]Sensor, 0)
	for i := 1; i <= numberOfSensors; i++ {
		number := strconv.Itoa(i)
		if i < 10 {
			number = "0" + number
		}
		sensors = append(sensors, Sensor{ID: "Sim-Sensor-" + number, Timestamp: time.Now(), Location: centralPoint})
	}
	return &SensorManager{api: api, sensors: sensors}
}

func (m SensorManager) Update() error {
	for _, s := range m.sensors {
		s.Timestamp = time.Now()
		latDelta := (float64(rand.Intn(7)-3) * .0001)
		lonDelta := (float64(rand.Intn(7)-3) * .0001)
		s.Location.Latitude = math.Round((s.Location.Latitude+latDelta)*100000) / 100000
		s.Location.Longitude = math.Round((s.Location.Longitude+lonDelta)*100000) / 100000
		if err := m.api.SendReading(s); err != nil {
			return err
		}
		log.Print("Updated sensor:", s, "DELTA:", latDelta, lonDelta)
	}
	return nil
}
