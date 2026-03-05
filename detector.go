package main

import (
	"time"
)

type Detector struct {
	ID        string           `json:"deploymentID"`
	Start     time.Time        `json:"start"`
	End       time.Time        `json:"end"`
	Type      string           `json:"sensorType"`
	Location  Point            `json:"location"`
	SpeciesID string           `json:"speciesID"`
	Animals   []DetectedAnimal `json:"animals"`
}

type DetectedAnimal struct {
	Behaviour   string `json:"behaviour"`
	Condition   string `json:"condition"`
	Confidence  int    `json:"confidence"`
	Description string `json:"description"`
	LifeStage   string `json:"lifeStage"`
	Sex         string `json:"sex"`
}
