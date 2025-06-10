package main

import "image"

type Target struct {
	Costumes       []image.Image
	CurrentCostume int64
	Events         []Event
}

type Costume struct {
	Target        Target
	X             float64
	Y             float64
	Directions    float64
	Size          float64
	Draggable     bool
	RotationStyle string
	Visible       bool
}
