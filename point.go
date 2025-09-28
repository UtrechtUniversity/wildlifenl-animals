package main

import "strconv"

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (p Point) Location() string {
	return strconv.FormatFloat(p.Latitude, 'f', -1, 64) + "," + strconv.FormatFloat(p.Longitude, 'f', -1, 64)
}
