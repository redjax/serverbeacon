package validators

import "fmt"

// Validate port is a digit 1-65535, return error if invalid
func ValidateTCPPort(p int64) error {
	if !isValidTCPPort(p) {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", p)
	}

	return nil
}

// Return true if port is valid
func isValidTCPPort(p int64) bool {
	if p < 1 || p > 65535 {
		fmt.Printf("invalid port: %d (must be 1-65535)", p)

		return false
	}

	return true
}
