package main

import (
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

type SpeedData struct {
	Download float64
	Upload   float64
}

func GetNetworkSpeed() (SpeedData, error) {
	startCounters, err := net.IOCounters(false)
	if err != nil {
		return SpeedData{}, err
	}
	if len(startCounters) == 0 {
		return SpeedData{}, nil
	}

	time.Sleep(1 * time.Second)

	endCounters, err := net.IOCounters(false)
	if err != nil {
		return SpeedData{}, err
	}
	if len(endCounters) == 0 {
		return SpeedData{}, nil
	}

	return SpeedData{
		Download: float64(endCounters[0].BytesRecv - startCounters[0].BytesRecv),
		Upload:   float64(endCounters[0].BytesSent - startCounters[0].BytesSent),
	}, nil
}
