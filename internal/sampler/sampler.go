package sampler

import (
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

// NetStats defines an interface for retrieving network I/O counters.
type NetStats interface {
	GetTotalBytes() (sent, recv uint64)
}

// GopsutilNetStats implements NetStats using gopsutil.
type GopsutilNetStats struct{}

func (g *GopsutilNetStats) GetTotalBytes() (sent, recv uint64) {
	counters, err := net.IOCounters(false)
	if err != nil || len(counters) == 0 {
		return 0, 0
	}

	for _, c := range counters {
		sent += c.BytesSent
		recv += c.BytesRecv
	}
	return sent, recv
}

// SampleResult holds the calculated rates.
type SampleResult struct {
	Sent  float64 `json:"sent"`
	Recv  float64 `json:"recv"`
	Total float64 `json:"total"`
}

// Sampler tracks network traffic and calculates rates.
type Sampler struct {
	stats      NetStats
	prevSent   uint64
	prevRecv   uint64
	prevTime   time.Time
	now        func() time.Time
	lastResult SampleResult
}

// NewSampler creates a new Sampler instance.
func NewSampler(stats NetStats) *Sampler {
	return &Sampler{
		stats: stats,
		now:   time.Now,
	}
}

// Sample calculates the current bytes per second.
func (s *Sampler) Sample() SampleResult {
	currentSent, currentRecv := s.stats.GetTotalBytes()
	now := s.now()

	if s.prevTime.IsZero() {
		s.prevSent = currentSent
		s.prevRecv = currentRecv
		s.prevTime = now
		s.lastResult = SampleResult{}
		return s.lastResult
	}

	timeDiff := now.Sub(s.prevTime).Seconds()
	if timeDiff <= 0 {
		return s.lastResult
	}

	if currentSent < s.prevSent || currentRecv < s.prevRecv {
		s.prevSent = currentSent
		s.prevRecv = currentRecv
		s.prevTime = now
		return s.lastResult
	}

	sentDiff := float64(currentSent - s.prevSent)
	recvDiff := float64(currentRecv - s.prevRecv)

	s.lastResult = SampleResult{
		Sent:  sentDiff / timeDiff,
		Recv:  recvDiff / timeDiff,
		Total: (sentDiff + recvDiff) / timeDiff,
	}

	s.prevSent = currentSent
	s.prevRecv = currentRecv
	s.prevTime = now

	return s.lastResult
}

// ToJSON returns the last sampled rates as a JSON string.
func (s *Sampler) ToJSON() string {
	b, err := json.Marshal(s.lastResult)
	if err != nil {
		return `{"sent":0,"recv":0,"total":0}`
	}
	return string(b)
}
