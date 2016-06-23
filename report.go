package main

import (
	"time"
)

// ResourceUsage is a wrapped result of getrusage(2)
type ResourceUsage struct {
	Real   time.Duration `json:"real"`
	User   time.Duration `json:"user"`
	System time.Duration `json:"system"`
}

// Report is a JSON-compatible report for ntimes
type Report struct {
	Average ResourceUsage `json:"average"`
	Max     ResourceUsage `json:"max"`
	Min     ResourceUsage `json:"min"`
	// Flaky is the percentage of failures (0.0-100.0).
	Flaky float64 `json:"flaky"`
}
