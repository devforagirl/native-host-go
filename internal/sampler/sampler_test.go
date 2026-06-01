package sampler

import (
	"encoding/json"
	"testing"
	"time"
)

type MockNetStatsProvider struct {
	sent uint64
	recv uint64
}

func (m *MockNetStatsProvider) GetTotalBytes() (sent, recv uint64) {
	return m.sent, m.recv
}

func almostEqual(a, b float64) bool {
	return (a-b) < 1e-9 && (b-a) < 1e-9
}

func TestSampler_Sample(t *testing.T) {
	startTime := time.Now()

	tests := []struct {
		name     string
		initial  [2]uint64 // sent, recv
		after    [2]uint64 // sent, recv
		duration time.Duration
		wantSent float64
		wantRecv float64
	}{
		{
			name:     "first sample should be 0",
			initial:  [2]uint64{100, 200},
			after:    [2]uint64{100, 200},
			duration: time.Second,
			wantSent: 0,
			wantRecv: 0,
		},
		{
			name:     "calculate rate correctly",
			initial:  [2]uint64{100, 200},
			after:    [2]uint64{1100, 1200}, // +1000 sent, +1000 recv
			duration: time.Second,
			wantSent: 1000,
			wantRecv: 1000,
		},
		{
			name:     "calculate rate with half second",
			initial:  [2]uint64{100, 200},
			after:    [2]uint64{600, 700}, // +500 sent, +500 recv
			duration: 500 * time.Millisecond,
			wantSent: 1000,
			wantRecv: 1000,
		},
		{
			name:     "counter reset protection",
			initial:  [2]uint64{1000, 1000},
			after:    [2]uint64{100, 100}, // reset
			duration: time.Second,
			wantSent: 0,
			wantRecv: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockNetStatsProvider{sent: tt.initial[0], recv: tt.initial[1]}
			s := NewSampler(mock)

			// Mock time
			currentTime := startTime
			s.now = func() time.Time { return currentTime }

			// First sample to initialize prev values
			s.Sample()

			// Update mock values and time
			currentTime = currentTime.Add(tt.duration)
			mock.sent = tt.after[0]
			mock.recv = tt.after[1]

			res := s.Sample()

			if !almostEqual(res.Sent, tt.wantSent) {
				t.Errorf("Sample().Sent = %v, want %v", res.Sent, tt.wantSent)
			}
			if !almostEqual(res.Recv, tt.wantRecv) {
				t.Errorf("Sample().Recv = %v, want %v", res.Recv, tt.wantRecv)
			}
		})
	}
}

func TestSampler_ToJSON(t *testing.T) {
	mock := &MockNetStatsProvider{sent: 100, recv: 200}
	s := NewSampler(mock)

	// We don't need to sample for ToJSON if it just returns current state,
	// but the requirement says it returns sent, recv, total.
	// Usually this refers to the last sampled rates.

	s.now = func() time.Time { return time.Now() }
	s.Sample() // Init

	// Mock some rates
	s.prevSent = 0
	s.prevRecv = 0
	s.prevTime = time.Now().Add(-time.Second)

	mock.sent = 1000
	mock.recv = 2000

	s.Sample()

	jsonStr := s.ToJSON()
	var result map[string]float64
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !almostEqual(result["sent"], 1000) {
		t.Errorf("JSON sent = %v, want 1000", result["sent"])
	}
	if !almostEqual(result["recv"], 2000) {
		t.Errorf("JSON recv = %v, want 2000", result["recv"])
	}
	if !almostEqual(result["total"], 3000) {
		t.Errorf("JSON total = %v, want 3000", result["total"])
	}
}
