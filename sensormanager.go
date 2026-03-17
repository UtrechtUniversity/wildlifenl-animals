package main

import (
	"math/rand"
	"strconv"
	"time"
)

type SensorManager struct {
	api     *WildlifeNLAPI
	sensors []Sensor
}

func NewSensorManager(centroid Point, numberOfSensors int, api *WildlifeNLAPI) *SensorManager {
	sensors := make([]Sensor, 0)
	for i := 1; i <= numberOfSensors; i++ {
		number := strconv.Itoa(i)
		if i < 10 {
			number = "0" + number
		}
		sensors = append(sensors, Sensor{ID: "Sim-Sensor-" + number, Timestamp: time.Now(), Location: centroid})
	}
	return &SensorManager{api: api, sensors: sensors}
}

func (m SensorManager) Update() error {
	for _, s := range m.sensors {
		now := time.Now()
		if now.Hour() > 21 || now.Hour() < 7 {
			continue
		}
		s.Timestamp = now
		step := rand.Intn(22) + 10
		cor := int(step / 2)
		s.Location.Latitude += (float64(rand.Intn(step)-cor) * .00017)
		s.Location.Longitude += (float64(rand.Intn(step)-cor) * .00017)
		if err := m.api.SendReading(s); err != nil {
			return err
		}
	}
	return nil
}
