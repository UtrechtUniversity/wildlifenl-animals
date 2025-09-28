package main

import "time"

type Sensor struct {
	ID        string    `json:"sensorID"`
	Timestamp time.Time `json:"timestamp"`
	Location  Point     `json:"location"`
}
