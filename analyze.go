package main

import (
	"image"
	"image/color"
)

//Segment is a type that stores a section of the row which the bot can just hold down the mouse button
type Segment struct {
	Start int
	End   int
}

type ConditionOptions struct {
	Pass      rune //? Which pass to compare to "RGBA", pick one letter
	Threshold int  //? The threshold to distinguish between mouse down and up
	Inverse   bool //? Brighter than threshold = Mouse up
}

//Analyze constructs a list of segments in rows to optimize image raster
func Analyze(img image.Image, options ConditionOptions) [][]Segment {
	rowsegments := [][]Segment{}
	height := img.Bounds().Max.Y
	for i := 0; i < height; i++ {
		rowsegments = append(rowsegments, AnalyzeRow(img, i, options))
	}

	return rowsegments
}

func AnalyzeRow(img image.Image, row_num int, options ConditionOptions) []Segment {
	segments := []Segment{}
	rowlength := img.Bounds().Max.X

	currentStart := 0
	for i := 0; i < rowlength; i++ {
		if AnalyzeCondition(img.At(i, row_num), options) {
			currentStart = i
		} else {
			segments = append(segments, Segment{currentStart, i})
		}
	}

	return segments
}

//AnalyzeCondition returns true to trigger mouse down
func AnalyzeCondition(color color.Color, options ConditionOptions) bool {
	r, g, b, a := color.RGBA()

	switch options.Pass {
	case 'R':
		if r > uint32(options.Threshold) {
			if !options.Inverse {
				return true
			}
		}
	case 'G':
		if g > uint32(options.Threshold) {
			if !options.Inverse {
				return true
			}
		}
	case 'B':
		if b > uint32(options.Threshold) {
			if !options.Inverse {
				return true
			}
		}
	case 'A':
		if a > uint32(options.Threshold) {
			if !options.Inverse {
				return true
			}
		}
	}

	return false
}
