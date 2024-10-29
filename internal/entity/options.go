package entity

import (
	"io"
	"time"
)

// ReporterOptions represents the options for the Reporter.
type ReporterOptions struct {
	// Since is the time to let the Reporter know when to start reporting.
	Since *time.Time
	// Writer is the io.Writer to write the report to.
	Writer io.Writer
}

// ReporterOption represents a functional option for the Reporter.
type ReporterOption func(*ReporterOptions)

// WithSince returns a ReporterOption that sets the Since field.
func WithSince(t time.Time) func(o *ReporterOptions) {
	return func(o *ReporterOptions) {
		o.Since = &t
	}
}

// WithWriter returns a ReporterOption that sets the Writer field.
func WithWriter(w io.Writer) func(o *ReporterOptions) {
	return func(o *ReporterOptions) {
		o.Writer = w
	}
}
