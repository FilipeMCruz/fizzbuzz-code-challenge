package services

import (
	"errors"
	"sync"
)

var ErrNoRequestsReceived = errors.New("no requests received")

// Stats is a simple struct that stores info about the requests made and also
// the most frequent one.
//
// A mutex is used to ensure that the request info returned is always the most frequent one,
// in exchange for a slight worst performance
type Stats struct {
	mostFrequent        map[string]int
	mostFrequentRequest string
	mostFrequentCount   int

	m sync.Mutex
}

// NewStats creates a Stats struct that is feed requests info from ch
func NewStats(ch <-chan string) *Stats {
	s := &Stats{
		mostFrequent:        make(map[string]int),
		mostFrequentRequest: "",
		mostFrequentCount:   0,
	}

	go s.run(ch)

	return s
}

func (s *Stats) run(ch <-chan string) {
	for req := range ch {
		s.increment(req)
	}
}

// increment increases the counter for the request received while keeping the most frequent one updated
func (s *Stats) increment(req string) {
	s.m.Lock()
	defer s.m.Unlock()

	count := s.mostFrequent[req]
	s.mostFrequent[req] = count + 1

	if count+1 > s.mostFrequentCount {
		s.mostFrequentRequest = req
		s.mostFrequentCount = count + 1
	}
}

// MostFrequent returns the most frequent request or an error if no request is found
func (s *Stats) MostFrequent() (string, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.mostFrequentCount == 0 {
		return "", ErrNoRequestsReceived
	}

	return s.mostFrequentRequest, nil
}
