package main

import (
	"image"

	"github.com/go-vgo/robotgo"
)

func draw(img image.Image, startX, startY int) {
	drawing = true

	draw_segments := Analyze(img, ConditionOptions{'R', 100, false})

	for y, row := range draw_segments {
		columnLevel := startY + y
		for _, segment := range row {
			robotgo.MoveMouse(segment.Start+startX, columnLevel)
			robotgo.Drag(segment.End+startX, columnLevel)
		}
	}

	robotgo.EventEnd()
}
