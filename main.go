package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"flowmeter-host/internal/messenger"
	"flowmeter-host/internal/registrar"
	"flowmeter-host/internal/sampler"
)

const hostName = "com.flowmeter.host"

func main() {
	register := flag.Bool("register", false, "Register the host with the browser")
	flag.Parse()

	if *register {
		binaryPath, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get executable path: %v", err)
		}

		reg := registrar.NewRegistrar()
		if err := reg.Register(hostName, binaryPath); err != nil {
			log.Fatalf("Registration failed: %v", err)
		}

		fmt.Println("Successfully registered FlowMeter Native Host.")
		os.Exit(0)
	}

	// Main Host Mode
	s := sampler.NewSampler(&sampler.GopsutilNetStats{})

	// We use a ticker to sample every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// The browser connects via stdin/stdout.
	// We send data periodically.
	for {
		select {
		case <-ticker.C:
			// 1. Sample current rates
			s.Sample()

			// 2. Convert to JSON
			jsonMsg := s.ToJSON()

			// 3. Send via messenger protocol to stdout
			if err := messenger.SendMessage(os.Stdout, []byte(jsonMsg)); err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}
		}
	}
}
