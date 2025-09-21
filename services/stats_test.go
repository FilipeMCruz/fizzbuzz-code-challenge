package services_test

import (
	"errors"
	"fizzbuzz-code-challenge/services"
	"reflect"
	"testing"
)

func TestStats(t *testing.T) {
	type testCase struct {
		description string
		input       []string
		ret         string
		err         error
	}

	testCases := []testCase{
		{
			description: "no requests",
			err:         services.ErrNoRequestsReceived,
		},
		{
			description: "one request",
			input:       []string{"single"},
			ret:         "single",
		},
		{
			description: "two requests",
			input:       []string{"single", "other"},
			ret:         "single",
		},
		{
			description: "three requests",
			input:       []string{"single", "other", "other"},
			ret:         "other",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ch := make(chan string)
			defer close(ch)

			s := services.NewStats()
			for i := range tc.input {
				s.Increment(tc.input[i])
			}

			ret, err := s.MostFrequent()

			if !errors.Is(tc.err, err) {
				t.Errorf("got %v, expected %v", err, tc.err)
			}

			if !reflect.DeepEqual(ret, tc.ret) {
				t.Errorf("got %v, expected %v", ret, tc.ret)
			}
		})
	}
}
