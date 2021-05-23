package lib

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
	Pass      rune //? Which pass to compare to "RGBA", pick one letter (if it is 'V' then it will use HSV value comparison)
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

	currentStart := -1 //? Is -1 while hasn't set a start point for the segment
	for i := 0; i < rowlength; i++ {
		if AnalyzeCondition(img.At(i, row_num), options) {
			if currentStart < 0 {
				currentStart = i
			}
		} else {
			if currentStart > -1 {
				segments = append(segments, Segment{currentStart, i})
				currentStart = -1
			}
		}
	}

	if currentStart > -1 {
		segments = append(segments, Segment{currentStart, rowlength})
	}

	return segments
}

//AnalyzeCondition returns true to trigger mouse down
func AnalyzeCondition(color color.Color, options ConditionOptions) bool {
	r, g, b, a := color.RGBA()
	r = r / 255
	g = g / 255
	b = b / 255
	a = a / 255

	switch options.Pass {
	case 'R':
		if r > uint32(options.Threshold) {
			return !options.Inverse
		}
	case 'G':
		if g > uint32(options.Threshold) {
			return !options.Inverse
		}
	case 'B':
		if b > uint32(options.Threshold) {
			return !options.Inverse
		}
	case 'A':
		if a > uint32(options.Threshold) {
			return !options.Inverse
		}
	case 'V':
		if max(int(r), int(g), int(b)) > options.Threshold {
			return !options.Inverse
		}
	}

	return options.Inverse
}
