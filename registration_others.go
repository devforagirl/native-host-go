//go:build !windows && !darwin

package main

import "fmt"

func doPlatformRegistration(manifestBytes []byte) error {
	return fmt.Errorf("unsupported operating system")
}
