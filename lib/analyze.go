package lib

import (
	"image"
)

//Segment is a type that stores a section of the row which the bot can just hold down the mouse button
type Segment struct {
	Start int
	End   int
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

//AnalyzeRow constructs segments from a row
func AnalyzeRow(img image.Image, row_num int, options ConditionOptions) []Segment {
	segments := []Segment{}
	rowlength := img.Bounds().Max.X

	currentStart := -1 //? Is -1 while hasn't set a start point for the segment
	for i := 0; i < rowlength; i++ {
		if ColorCondition(img.At(i, row_num), options) {
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
