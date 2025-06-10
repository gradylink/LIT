package main

import (
	"image"
	"math"
)

type Sprite struct {
	Costumes       []image.Image
	CurrentCostume int64
	Events         []Event
	X              float64
	Y              float64
	Direction      float64
	Size           float64
	Draggable      bool
	RotationStyle  int8 // 0 is All Around, 1 is Left-Right, and 2 is Don't Rotate
	Visible        bool
}

func (s *Sprite) Run() {
	for _, event := range s.Events {
		if event.Event == "event_whenflagclicked" {
			go event.Method(s.Render, nil, nil)
		}
	}
}

func (s *Sprite) Render() {}

func (s *Sprite) Motion_Move(dis float64) {
	s.X += dis * math.Cos((90-s.Direction)*180/math.Pi)
	s.Y += dis * math.Sin((90-s.Direction)*180/math.Pi)
}
