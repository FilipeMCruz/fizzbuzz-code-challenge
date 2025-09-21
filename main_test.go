package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"testing"
)

type req struct {
	method       string
	url          string
	expectedBody string
	expectedCode int
}

type testCase struct {
	description string
	reqs        []req
	err         error
}

func TestRun(t *testing.T) {
	testCases := []testCase{
		{
			description: "stats -> fizzbuzz 1 -> stats -> fizzbuzz 2 -> stats -> fizzbuzz 2 -> stats",
			reqs: []req{
				{
					method:       "GET",
					url:          "/api/v1/stats",
					expectedCode: http.StatusNotFound,
					expectedBody: `{"error":"no requests received"}`,
				},
				{
					method:       "GET",
					url:          "/api/v1/fizzbuzz?int1=3",
					expectedCode: http.StatusBadRequest,
					expectedBody: `{"error":"invalid query param: int2"}`,
				},
				{
					method:       "GET",
					url:          "/api/v1/stats",
					expectedCode: http.StatusOK,
					expectedBody: `{"most_frequent":"/api/v1/stats"}`,
				},
				{
					method:       "GET",
					url:          "/api/v1/fizzbuzz?int1=3&int2=5&limit=10&str1=a&str2=b",
					expectedCode: http.StatusOK,
					expectedBody: `{"values":["1","2","a","4","b","a","7","8","a","b"],"total":10}`,
				},
				{
					method:       "GET",
					url:          "/api/v1/fizzbuzz?int1=3&int2=5&limit=10&str1=a&str2=b",
					expectedCode: http.StatusOK,
					expectedBody: `{"values":["1","2","a","4","b","a","7","8","a","b"],"total":10}`,
				},
				{
					method:       "GET",
					url:          "/api/v1/fizzbuzz?int1=3&int2=5&limit=10&str1=a&str2=b",
					expectedCode: http.StatusOK,
					expectedBody: `{"values":["1","2","a","4","b","a","7","8","a","b"],"total":10}`,
				},
				{
					method:       "GET",
					url:          "/api/v1/stats",
					expectedCode: http.StatusOK,
					expectedBody: `{"most_frequent":"/api/v1/fizzbuzz?int1=3\u0026int2=5\u0026limit=10\u0026str1=a\u0026str2=b"}`,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx, stop := context.WithCancel(context.Background())
			defer stop()

			port, err := getFreePort()
			if err != nil {
				t.Fatal(err)
			}

			noti := make(chan struct{})

			running := func() {
				noti <- struct{}{}
			}

			go func() {
				startErr := start(ctx, stop, running, port)

				if !errors.Is(tc.err, startErr) {
					t.Errorf("got %v, expected %v", startErr, tc.err)
				}
			}()

			<-noti

			runTest(t, tc, port)
		})
	}
}

func runTest(t *testing.T, tc testCase, port int) {
	for _, req := range tc.reqs {
		r, _ := http.NewRequest(req.method, "http://localhost:"+strconv.Itoa(port)+req.url, nil)

		resp, reqErr := http.DefaultClient.Do(r)
		if reqErr != nil {
			t.Fatal(reqErr)
		}

		if resp.StatusCode != req.expectedCode {
			t.Errorf("got %d, expected %d", resp.StatusCode, req.expectedCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if string(body) != req.expectedBody {
			t.Errorf("got %s, expected %s", string(body), req.expectedBody)
		}
	}
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return -1, err
	}

	l, tcpErr := net.ListenTCP("tcp", addr)
	if tcpErr != nil {
		return -1, err
	}

	defer func(l *net.TCPListener) {
		_ = l.Close()
	}(l)
	return l.Addr().(*net.TCPAddr).Port, nil
}
