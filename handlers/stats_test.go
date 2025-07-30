package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildStatsHandler(t *testing.T) {
	type testCase struct {
		description        string
		input              []string
		url                string
		method             string
		expectedResponse   []byte
		expectedStatusCode int
	}

	testCases := []testCase{
		{
			description:        "no requests made so far",
			url:                "/stats",
			method:             "GET",
			expectedResponse:   []byte(`{"error":"no requests received"}`),
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description:        "one request",
			input:              []string{"single"},
			url:                "/stats",
			method:             "GET",
			expectedResponse:   []byte(`{"most_frequent":"single"}`),
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "three requests",
			input:              []string{"single", "other", "other"},
			url:                "/stats",
			method:             "GET",
			expectedResponse:   []byte(`{"most_frequent":"other"}`),
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ch := make(chan string)
			defer close(ch)

			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			h := BuildStatsHandler(ch)

			for i := range tc.input {
				ch <- tc.input[i]
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
