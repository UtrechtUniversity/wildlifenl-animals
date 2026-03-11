package main

import (
	"math/rand/v2"
	"time"
)

type DetectionManager struct {
	api       *WildlifeNLAPI
	species   []Species
	detectors []Detector
}

func NewDetectionManager(centroid Point, species []Species, api *WildlifeNLAPI) *DetectionManager {

	location := centroid
	acoustic1 := Detector{ID: "Sim-Detector-01", Type: "acoustic", Location: location}

	location.Longitude -= 0.0005
	cameraTrap1 := Detector{ID: "Sim-Detector-02", Type: "visual", Location: location}

	location.Longitude += 0.001
	location.Latitude += 0.0005
	cameraTrap2 := Detector{ID: "Sim-Detector-03", Type: "visual", Location: location}

	location.Latitude -= 0.001
	cameraTrap3 := Detector{ID: "Sim-Detector-04", Type: "visual", Location: location}

	return &DetectionManager{api: api, species: species, detectors: []Detector{acoustic1, cameraTrap1, cameraTrap2, cameraTrap3}}
}

func (m DetectionManager) Update() error {
	for _, d := range m.detectors {
		if rand.IntN(12) != 7 { // Lucky number ;-)
			continue
		}
		d.End = time.Now()
		d.Start = time.Now().Local().Add(3 * time.Minute * -1)
		d.SpeciesID = m.species[rand.IntN(len(m.species))].ID
		d.Animals = m.generateRandomDetectedAnimals()
		if err := m.api.SendDetection(d); err != nil {
			return err
		}
	}
	return nil
}

func (m DetectionManager) generateRandomDetectedAnimals() []DetectedAnimal {
	min := 1
	max := 5
	conditions := []string{"healthy", "impaired", "dead"}
	lifeStages := []string{"infant", "adolescent", "adult"}
	sexes := []string{"female", "male"}
	animals := make([]DetectedAnimal, 0)
	for range rand.IntN(max-min) + min {
		a := DetectedAnimal{
			Behaviour:   "Present",
			Condition:   conditions[rand.IntN(len(conditions))],
			Confidence:  rand.IntN(100),
			Description: "Animal detected.",
			LifeStage:   lifeStages[rand.IntN(len(lifeStages))],
			Sex:         sexes[rand.IntN(len(sexes))],
		}
		animals = append(animals, a)
	}
	return animals
}
