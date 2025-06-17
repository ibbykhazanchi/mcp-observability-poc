package main

import (
	"mcp-observability/handlers"

	server "github.com/ckanthony/gin-mcp"

	"github.com/gin-gonic/gin"
)

type GetSpansParams struct {
	SpanID     string            `form:"spanId" json:"spanId,omitempty" jsonschema:"description=Span ID"`
	TraceID    string            `form:"traceId" json:"traceId,omitempty" jsonschema:"description=Trace ID"`
	StartTime  string            `form:"startTime" json:"startTime,omitempty" jsonschema:"description=Start Time"`
	EndTime    string            `form:"endTime" json:"endTime,omitempty" jsonschema:"description=End Time"`
	Attributes map[string]string `form:"attributes" json:"attributes,omitempty" jsonschema:"description=Attributes"`
	Status     string            `form:"status" json:"status,omitempty" jsonschema:"description=Status"`
	Events     []string          `form:"events" json:"events,omitempty" jsonschema:"description=Events"`
	Links      []string          `form:"links" json:"links,omitempty" jsonschema:"description=Links"`
	Kind       string            `form:"kind" json:"kind,omitempty" jsonschema:"description=Kind"`
	Name       string            `form:"name" json:"name,omitempty" jsonschema:"description=Name"`
}

func main() {
	r := gin.Default()

	// Telemetry routes
	r.GET("/spans", handlers.GetSpans)

	mcp := server.New(r, &server.Config{
		Name:        "mcp-observability",
		Description: "Observability MCP server",
		BaseURL:     "http://localhost:8080",
	})

	mcp.RegisterSchema("GET", "/spans", GetSpansParams{}, nil)

	mcp.Mount("/mcp")

	// Start the serve
	r.Run(":8080")
}
