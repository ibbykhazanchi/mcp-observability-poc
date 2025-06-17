package models

import "time"

// Span represents a single operation within a trace
type Span struct {
	TraceID    string                 `json:"traceId"`
	SpanID     string                 `json:"spanId"`
	Name       string                 `json:"name"`
	Kind       string                 `json:"kind"`
	StartTime  time.Time              `json:"startTime"`
	EndTime    time.Time              `json:"endTime"`
	Attributes map[string]interface{} `json:"attributes"`
	Status     Status                 `json:"status"`
	Events     []Event                `json:"events"`
	Links      []Link                 `json:"links"`
}

// Status represents the status of a span
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Event represents an event that occurred during a span
type Event struct {
	Time       time.Time              `json:"time"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Link represents a link to another span
type Link struct {
	TraceID    string                 `json:"traceId"`
	SpanID     string                 `json:"spanId"`
	Attributes map[string]interface{} `json:"attributes"`
} 