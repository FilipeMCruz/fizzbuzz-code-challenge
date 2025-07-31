package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
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

			port, err := GetFreePort()
			if err != nil {
				t.Fatal(err)
			}

			go func() {
				err := start(ctx, stop, port)

				if !reflect.DeepEqual(tc.err, err) {
					t.Errorf("got %v, expected %v", err, tc.err)
				}
			}()

			time.Sleep(time.Second)

			for _, req := range tc.reqs {
				r, _ := http.NewRequest(req.method, "http://localhost:"+strconv.Itoa(port)+req.url, nil)

				resp, err := http.DefaultClient.Do(r)
				if err != nil {
					t.Fatal(err)
				}

				if resp.StatusCode != req.expectedCode {
					t.Errorf("got %d, expected %d", resp.StatusCode, req.expectedCode)
				}

				body, _ := io.ReadAll(resp.Body)
				if string(body) != req.expectedBody {
					t.Errorf("got %s, expected %s", string(body), req.expectedBody)
				}
			}
		})
	}
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err == nil {
		l, err := net.ListenTCP("tcp", addr)
		if err == nil {
			defer func(l *net.TCPListener) {
				_ = l.Close()
			}(l)
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}

	return -1, err
}
