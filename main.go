package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	fileInfo, _ := os.Stdin.Stat()
	isTty := (fileInfo.Mode() & os.ModeCharDevice) != 0

	if isTty {
		HandleRegistration()
	} else {
		handleService()
	}
}

func handleService() {
	fmt.Fprintln(os.Stderr, "Mode: Service (Non-TTY)")

	for {
		speed, err := GetNetworkSpeed()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting network speed: %v\n", err)
			continue
		}

		payload := map[string]string{
			"type":     "speed_update",
			"download": formatBytes(speed.Download),
			"upload":   formatBytes(speed.Upload),
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			continue
		}

		err = WriteMessage(os.Stdout, jsonPayload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing message: %v\n", err)
			os.Exit(1)
		}
	}
}

func formatBytes(bytes float64) string {
	units := []string{"B/s", "KB/s", "MB/s", "GB/s", "TB/s", "PB/s", "EB/s"}
	var i int
	for bytes >= 1024 && i < len(units)-1 {
		bytes /= 1024
		i++
	}
	return fmt.Sprintf("%.2f %s", bytes, units[i])
}
