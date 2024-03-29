package effects

import (
	"math"
)

// moduloFloat64 implements math.Mod with proper negative number support
func moduloFloat64(x, y float64) float64 {
	for x < 0 {
		x += y
	}
	return math.Mod(x, y)
}

// moduloInt implements % with proper negative number support
func moduloInt[K int | int64](x, y K) K {
	for x < 0 {
		x += y
	}
	return x % y
}

// getDirectionFactor returns a factor that can be multiplied with a position offset so that is either reversed
// (negative direction) or not
func getDirectionFactor(reversed bool) int {
	if reversed {
		return -1
	}
	return 1
}

// getNextPosition calculates the next position of a point on the LED strip for constant movements
// It updates the position parameter and returns the current position as integer.
// We calculate based on floats to have a better precision.
func getNextPosition(position *float64, ledsPerSecond float64, numberLeds int, nanoseconds int64, reversed bool) int {
	*position = *position + float64(getDirectionFactor(reversed))*((ledsPerSecond*float64(nanoseconds))/1000000000.0)

	// normalize to 0 <= pos < number_leds
	*position = moduloFloat64(*position, float64(numberLeds))

	// return the current led as int
	return int(*position)
}

// mapPercent returns a value between min and max the corresponds to the specified percentage
func mapPercent[K int | int64 | float64](min K, max K, percent int) K {
	return min + ((max - min) * K(percent) / 100)
}

// hueToRGB converts HSV to RGB
// H: 0-360, S: 0-100, V: 0-100
// For the HSV input, S and V are 255 and H is variable.
func hsv2rgb(h int, s int, v int) (r byte, g byte, b byte) {
	hTmp := float64(h) / 60
	sTmp := float64(s) / 100
	vTmp := float64(v) / 100
	hi := math.Mod(math.Floor(hTmp), 6)
	f := hTmp - math.Floor(hTmp)
	p := 255 * vTmp * (1 - sTmp)
	q := 255 * vTmp * (1 - (sTmp * f))
	t := 255 * vTmp * (1 - (sTmp * (1 - f)))
	vTmp *= 255

	var result [3]float64

	switch hi {
	case 0:
		result[0] = math.Round(vTmp)
		result[1] = math.Round(t)
		result[2] = math.Round(p)
	case 1:
		result[0] = math.Round(q)
		result[1] = math.Round(vTmp)
		result[2] = math.Round(p)
	case 2:
		result[0] = math.Round(p)
		result[1] = math.Round(vTmp)
		result[2] = math.Round(t)
	case 3:
		result[0] = math.Round(p)
		result[1] = math.Round(q)
		result[2] = math.Round(vTmp)
	case 4:
		result[0] = math.Round(t)
		result[1] = math.Round(p)
		result[2] = math.Round(vTmp)
	case 5:
		result[0] = math.Round(vTmp)
		result[1] = math.Round(p)
		result[2] = math.Round(q)
	}
	return byte(result[0]), byte(result[1]), byte(result[2])
}
