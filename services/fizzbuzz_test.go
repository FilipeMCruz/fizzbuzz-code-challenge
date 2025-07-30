package services

import (
	"reflect"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	type testCase struct {
		description       string
		int1, int2, limit int
		str1, str2        string
		ret               []string
		err               error
	}

	testCases := []testCase{
		{
			description: "basic example: FizzBuzz",
			limit:       20,
			int1:        3,
			int2:        5,
			str1:        "fizz",
			str2:        "buzz",
			ret:         []string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz", "16", "17", "fizz", "19", "buzz"},
		},
		{
			description: "bad upper limit",
			limit:       -1,
			int1:        3,
			int2:        5,
			str1:        "fizz",
			str2:        "buzz",
			ret:         []string{},
		},
		{
			description: "basic example with int2 being an int1 multiple",
			limit:       10,
			int1:        2,
			int2:        4,
			str1:        "a",
			str2:        "b",
			ret:         []string{"1", "a", "3", "b", "5", "a", "7", "ab", "9", "a"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ret, err := FizzBuzz(tc.int1, tc.int2, tc.limit, tc.str1, tc.str2)

			if !reflect.DeepEqual(tc.err, err) {
				t.Errorf("got %v, expected %v", err, tc.err)
			}

			if !reflect.DeepEqual(ret, tc.ret) {
				t.Errorf("got %v, expected %v", ret, tc.ret)
			}
		})
	}
}
