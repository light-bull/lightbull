package effects

import "math"

// getNextPosition calculates the next position of a point on the LED strip for constant movements
// It updates the position parameter and returns the current position as integer.
// We calculate based on floats to have a better precision.
func getNextPosition(position *float64, ledsPerSecond float64, numberLeds int, nanoseconds int64) int {
	*position = *position + ((ledsPerSecond * float64(nanoseconds)) / 1000000000.0)

	// normalize to 0 <= pos < number_leds
	*position = math.Mod(*position, float64(numberLeds))

	// return the current led as int
	return int(*position)
}

// mapPercent returns a value between min and max the corresponds to the specified percentage
func mapPercent[K int | int64 | float64](min K, max K, percent int) K {
	return min + ((max - min) * K(percent) / 100)
}
