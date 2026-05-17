package compute

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// SturgesCoeff calculates the number of bins using Sturges' rule.
// It returns ceil(1 + 3.322 * log10(n)) or an error if n is not positive.
func SturgesCoeff(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("n has to be positive, got = %d", n)
	}

	return int(math.Ceil(1.0 + 3.322*math.Log10(float64(n)))), nil
}

// GetSample returns a random sample of the specified size from the input data slice.
// It uses Fisher-Yates shuffle to randomly select elements. Returns an error if the sample size is not positive or exceeds the data length.
func GetSample[T any](data []T, size int) ([]T, error) {
	if size <= 0 {
		return nil, fmt.Errorf("sample size must be positive, got = %d", size)
	}

	if size > len(data) {
		return nil, fmt.Errorf("sample size is out of bounds, got = %d bound = %d", size, len(data))
	}

	sample := make([]T, len(data))
	copy(sample, data)

	for i := range size {
		j := i + rand.IntN(len(sample)-i)
		sample[i], sample[j] = sample[j], sample[i]
	}

	return sample[:size], nil
}
