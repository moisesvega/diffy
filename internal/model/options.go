package model

import "time"

// ReporterOptions represents the options for the Reporter.
type ReporterOptions struct {
	// Since is the time to let the Reporter know when to start reporting.
	Since time.Time
}

// ReporterOption represents a functional option for the Reporter.
type ReporterOption func(*ReporterOption)

// WithSince returns a ReporterOption that sets the Since field.
func WithSince(t time.Time) func(o *ReporterOptions) {
	return func(o *ReporterOptions) {
		o.Since = t
	}
}
