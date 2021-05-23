package lib

import "image/color"

type ConditionOptions struct {
	Pass      rune //? Which pass to compare to "RGBA", pick one letter (if it is 'V' then it will use HSV value comparison)
	Threshold int  //? The threshold to distinguish between mouse down and up
	Inverse   bool //? Less bright than threshold = Mouse down
}

//ColorCondition returns true to trigger mouse down if the condition is met
func ColorCondition(color color.Color, options ConditionOptions) bool {
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
