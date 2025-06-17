package handlers

import (
	"encoding/json"
	"fmt"
	"mcp-observability/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSpans handles retrieval of spans by any field using query parameters
func GetSpans(c *gin.Context) {
	queries := c.Request.URL.Query()

	data, err := os.ReadFile("data/spans.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read spans data"})
		return
	}

	var spans []models.Span
	if err := json.Unmarshal(data, &spans); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse spans data"})
		return
	}

	// If no query params, return all spans
	if len(queries) == 0 {
		c.JSON(http.StatusOK, spans)
		return
	}

	var result []models.Span
	for _, span := range spans {
		match := true
		for key, values := range queries {
			if len(values) == 0 {
				continue
			}
			if !spanFieldMatches(span, key, values[0]) {
				match = false
				break
			}
		}
		if match {
			result = append(result, span)
		}
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No spans found matching query"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// spanFieldMatches checks if a span's field (including nested/attributes) matches the query value
func spanFieldMatches(span models.Span, key, value string) bool {
	switch key {
	case "traceId":
		return span.TraceID == value
	case "spanId":
		return span.SpanID == value
	case "name":
		return span.Name == value
	case "kind":
		return span.Kind == value
	case "startTime":
		return span.StartTime.Format(time.RFC3339) == value
	case "endTime":
		return span.EndTime.Format(time.RFC3339) == value
	case "status.code":
		return intToString(span.Status.Code) == value
	case "status.message":
		return span.Status.Message == value
	default:
		// attributes.foo
		if len(key) > 10 && key[:10] == "attributes" {
			attrKey := key[11:]
			if v, ok := span.Attributes[attrKey]; ok {
				return interfaceToString(v) == value
			}
		}
		// events.name, events.attributes.foo, links.traceId, etc. (optional, can be extended)
	}
	return false
}

func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}

func interfaceToString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return fmt.Sprintf("%d", t)
	case float64:
		return fmt.Sprintf("%v", t)
	case bool:
		return fmt.Sprintf("%v", t)
	default:
		return ""
	}
}

// PostSpansQuery handles querying spans with flexible filters via POST
func PostSpansQuery(c *gin.Context) {
	type Filter struct {
		Field string      `json:"field"`
		Op    string      `json:"op"` // e.g., eq, neq, contains
		Value interface{} `json:"value"`
	}

	var filters []Filter
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter format"})
		return
	}

	data, err := os.ReadFile("data/spans.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read spans data"})
		return
	}

	var spans []models.Span
	if err := json.Unmarshal(data, &spans); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse spans data"})
		return
	}

	var result []models.Span
	for _, span := range spans {
		match := true
		for _, filter := range filters {
			if !applyFilter(span, filter) {
				match = false
				break
			}
		}
		if match {
			result = append(result, span)
		}
	}

	c.JSON(http.StatusOK, result)
}

// applyFilter applies a single filter to a span
func applyFilter(span models.Span, filter interface{}) bool {
	f, ok := filter.(map[string]interface{})
	if !ok {
		return false
	}
	field, _ := f["field"].(string)
	op, _ := f["op"].(string)
	value := f["value"]

	switch field {
	case "traceId":
		return compare(span.TraceID, op, value)
	case "spanId":
		return compare(span.SpanID, op, value)
	case "name":
		return compare(span.Name, op, value)
	case "kind":
		return compare(span.Kind, op, value)
	case "startTime":
		return compare(span.StartTime, op, value)
	case "endTime":
		return compare(span.EndTime, op, value)
	case "status.code":
		return compare(span.Status.Code, op, value)
	case "status.message":
		return compare(span.Status.Message, op, value)
	default:
		// Support for attributes, events, links
		if len(field) > 10 && field[:10] == "attributes" {
			attrKey := field[11:]
			if v, ok := span.Attributes[attrKey]; ok {
				return compare(v, op, value)
			}
		}
		// TODO: Add support for events and links if needed
	}
	return false
}

// compare compares two values based on the op
func compare(a interface{}, op string, b interface{}) bool {
	switch op {
	case "eq":
		return a == b
	case "neq":
		return a != b
	case "contains":
		as, ok := a.(string)
		bs, ok2 := b.(string)
		return ok && ok2 && (len(bs) == 0 || (len(as) >= len(bs) && contains(as, bs)))
	}
	return false
}

func contains(a, b string) bool {
	return len(b) == 0 || (len(a) >= len(b) && (a == b || (len(a) > len(b) && (a[:len(b)] == b || contains(a[1:], b)))))
}
