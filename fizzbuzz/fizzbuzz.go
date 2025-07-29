package fizzbuzz

import (
	"errors"
	"strconv"
)

var ErrZeroInt1Value = errors.New("int1 value can't be 0")
var ErrZeroInt2Value = errors.New("int2 value can't be 0")

// FizzBuzz function follows the semantics and naming convention provided by
// the code challenge.
// - int1: multiples of this variable should be replaced by str1
// - int2: multiples of this variable should be replaced by str2
// - limit: upper limit for the iteration
// - str1: string that will replace multiples of int1
// - str2: string that will replace multiples of int2
//
// Challenge doesn't reference any constrains but there's a need to at least check for 0's in int1 and in2
func FizzBuzz(int1, int2, limit int, str1, str2 string) ([]string, error) {
	if int1 == 0 {
		return nil, ErrZeroInt1Value
	}

	if int2 == 0 {
		return nil, ErrZeroInt2Value
	}

	if limit < 1 {
		return []string{}, nil
	}

	accumulator := make([]string, limit)

	bothInt := int1 * int2
	bothStr := str1 + str2
	if int1 > int2 {
		tmpInt := int2
		int2 = int1
		int1 = tmpInt

		tmpStr := str2
		str2 = str1
		str1 = tmpStr
	}

	for curr := 1; curr <= limit; curr++ {
		idx := curr - 1
		switch {
		case curr%bothInt == 0:
			accumulator[idx] = bothStr
		case curr%int2 == 0:
			accumulator[idx] = str2
		case curr%int1 == 0:
			accumulator[idx] = str1
		default:
			accumulator[idx] = strconv.Itoa(curr)
		}
	}

	return accumulator, nil
}
