package neurogo

import (
	"fmt"
	"math"
)

// mean computes the mean of a float32 slice.
//
// Parameters:
//  - data: the slice of float32 values. Must not be empty.
//
// Returns:
//  - float32: the mean
//  - error: an error if one occurred, e.g., the slice was empty
func mean(data[]float32) (float32, error) {
	if len(data) == 0 {
		err := fmt.Errorf("mean: empty slice")
		return 0.0, err
	}

	var sum float32 = 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float32(len(data)), nil
}


// max computes the maximum of a float32 slice.
//
// Parameters:
//  - data: the slice of float32 values. Must not be empty.
//
// Returns:
//  - float32: the maximum
//  - error: an error if one occurred, e.g., the slice was empty
func max(data[]float32) (float32, error) {
	if len(data) == 0 {
		err := fmt.Errorf("max: empty slice")
		return 0.0, err
	}

	var max float32 = - math.MaxFloat32
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max, nil
}

// min computes the minimum of a float32 slice.
//
// Parameters:
//  - data: the slice of float32 values. Must not be empty.
//
// Returns:
//  - float32: the minimum
//  - error: an error if one occurred, e.g., the slice was empty
func min(data[]float32) (float32, error) {
	if len(data) == 0 {
		err := fmt.Errorf("min: empty slice")
		return 0.0, err
	}

	var min float32 = math.MaxFloat32
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min, nil
}