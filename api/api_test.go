package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildFizzBuzzHandler(t *testing.T) {
	type testCase struct {
		description        string
		url                string
		method             string
		expectedResponse   []byte
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			description:        "basic example: FizzBuzz",
			url:                "/fizzbuzz?int1=3&int2=5&limit=20&str1=fizz&str2=buzz",
			method:             "GET",
			expectedResponse:   []byte(`["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16","17","fizz","19","buzz"]`),
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "invalid request: int1 needs to be a number",
			url:                "/fizzbuzz?int1=a&int2=5&limit=20&str1=fizz&str2=buzz",
			method:             "GET",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"invalid query param: int1"}`),
		},
		{
			description:        "invalid request: int2 needs to be a number",
			url:                "/fizzbuzz?int1=1&int2=b&limit=20&str1=fizz&str2=buzz",
			method:             "GET",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"invalid query param: int2"}`),
		},
		{
			description:        "invalid request: limit needs to be a number",
			url:                "/fizzbuzz?int1=1&int2=4&limit=g&str1=fizz&str2=buzz",
			method:             "GET",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"invalid query param: limit"}`),
		},
		{
			description:        "invalid request: str1 needs to be defined",
			url:                "/fizzbuzz?int1=1&int2=4&limit=10&&str2=buzz",
			method:             "GET",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"missing query param: str1"}`),
		},
		{
			description:        "invalid request: str2 needs to be defined",
			url:                "/fizzbuzz?int1=1&int2=4&limit=10&str1=fizz",
			method:             "GET",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   []byte(`{"error":"missing query param: str2"}`),
		},
	}

	h := BuildFizzBuzzHandler()

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			resp := rr.Result()

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			body, _ := io.ReadAll(resp.Body)

			if !bytes.Equal(tc.expectedResponse, body) {
				t.Errorf("got %v, expected %v", string(body), string(tc.expectedResponse))
			}

			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("got %v, expected %v", resp.StatusCode, tc.expectedStatusCode)
			}
		})
	}
}
