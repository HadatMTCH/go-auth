package utils

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

func NullIntScan(a *int) int {
	if a != nil {
		return *a
	}

	return 0
}

func NullFloatScan(a *float64) float64 {
	if a != nil {
		return *a
	}

	return 0.0
}

func NullFloat32Scan(a *float32) float32 {
	if a != nil {
		return *a
	}
	return 0
}

func ScanIntToNullValue(a int) *int {
	if a == 0 {
		return nil
	}

	return &a
}

func NullFloat64ScanFromNullableString(a *string) float64 {
	if a != nil {
		value, err := strconv.ParseFloat(*a, 64)
		if err != nil {
			return 0.0
		}
		return value
	}
	return 0.0
}

func CountTotalPage(total, perPage int) int {
	if (total % perPage) > 0 {
		return (total / perPage) + 1
	} else {
		return total / perPage
	}
}

func CommaSeparated(v float64) string {
	sign := ""

	// Min float64 can't be negated to a usable value, so it has to be special cased.
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}

	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = fmt.Sprintf("%.0f", math.Floor(math.Mod(v, 1000)))
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

func GenerateRandomNumber(length int) (uint64, error) {
	if length < 1 {
		return 0, fmt.Errorf("length must be at least 1")
	}
	if length > 20 {
		return 0, fmt.Errorf("length must be less than 19 to fit in int64")
	}

	// Calculate the range for the first digit (1-9)
	min := big.NewInt(1)
	max := big.NewInt(9)

	// Generate first digit (1-9)
	firstDigit, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, min))
	if err != nil {
		return 0, fmt.Errorf("failed to generate first digit: %v", err)
	}

	// Generate remaining digits
	result := firstDigit.Uint64()
	if length > 1 {
		// Calculate the range for remaining digits (0-9)
		max = big.NewInt(10)
		for i := 1; i < length; i++ {
			digit, err := rand.Int(rand.Reader, max)
			if err != nil {
				return 0, fmt.Errorf("failed to generate digit at position %d: %v", i, err)
			}
			result = result*10 + digit.Uint64()
		}
	}

	return result, nil
}
