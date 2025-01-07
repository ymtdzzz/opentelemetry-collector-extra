package newrelic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConnectionData(t *testing.T) {
	input := `
[
   {
      "pid":49993,
      "language":"go",
      "agent_version":"3.35.1",
      "host":"test-host",
      "settings":{},
      "app_name":[
         "go-agent-example"
      ],
      "high_security":false,
      "labels": {},
      "environment":[],
      "identifier":"go-agent-example",
      "utilization":{},
      "security_policies":{},
      "metadata":{},
      "event_harvest_config":{}
   }
]`

	want := &Connection{
		Pid:      49993,
		Language: "go",
		Version:  "3.35.1",
		Host:     "test-host",
		AppName: []string{
			"go-agent-example",
		},
	}

	got, err := ParseConnectionData([]byte(input))
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestParseSpanEventData(t *testing.T) {
	input := `
[
   "12345",
   {
      "reservoir_size":2000,
      "events_seen":10
   },
   [
      [
         {
            "type":"Span",
            "traceId":"a6da3dd78ae8d119d7d48420205b0dab",
            "guid":"a0af6f0ece5219eb",
            "parentId":"01fa4322cf7f0e86",
            "transactionId":"a6da3dd78ae8d119",
            "sampled":true,
            "priority":1.869206,
            "timestamp":1736146550635,
            "duration":0.464587875,
            "name":"External/example.com/http/GET",
            "category":"http",
            "component":"http",
            "span.kind":"client"
         },
         {},
         {
            "http.method":"GET",
            "http.url":"https://example.com"
         }
      ],
      [
         {
            "type":"Span",
            "traceId":"a6da3dd78ae8d119d7d48420205b0dab",
            "guid":"01fa4322cf7f0e86",
            "transactionId":"a6da3dd78ae8d119",
            "sampled":true,
            "priority":1.869206,
            "timestamp":1736146550635,
            "duration":0.465100584,
            "name":"WebTransaction/Go/GET /external",
            "category":"generic",
            "nr.entryPoint":true,
            "transaction.name":"WebTransaction/Go/GET /external"
         },
         {},
         {
            "request.headers.host":"localhost:8000",
            "code.filepath":"go-agent/v3/examples/server/main.go",
            "request.uri":"/external",
            "request.headers.accept":"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
            "code.namespace":"main",
            "request.method":"GET",
            "code.lineno":160,
            "code.function":"external"
         }
      ]
   ]
]`
	want := &SpanEvent{
		RunID: "12345",
		EventInfo: &EventInfo{
			ReservoirSize: 2000,
			EventsSeen:    10,
		},
		Spans: []*Span{
			{
				ParsedSpan: &ParsedSpan{
					Type:          "Span",
					TraceID:       "a6da3dd78ae8d119d7d48420205b0dab",
					SpanID:        "a0af6f0ece5219eb",
					ParentID:      "01fa4322cf7f0e86",
					TransactionID: "a6da3dd78ae8d119",
					Sampled:       true,
					Priority:      1.869206,
					Timestamp:     1736146550635,
					Duration:      0.464587875,
					Name:          "External/example.com/http/GET",
					Category:      "http",
					Component:     "http",
					SpanKind:      "client",
				},
				UserAttributes: map[string]any{},
				AgentAttributes: map[string]any{
					"http.method": "GET",
					"http.url":    "https://example.com",
				},
			},
			{
				ParsedSpan: &ParsedSpan{
					Type:            "Span",
					TraceID:         "a6da3dd78ae8d119d7d48420205b0dab",
					SpanID:          "01fa4322cf7f0e86",
					TransactionID:   "a6da3dd78ae8d119",
					Sampled:         true,
					Priority:        1.869206,
					Timestamp:       1736146550635,
					Duration:        0.465100584,
					Name:            "WebTransaction/Go/GET /external",
					Category:        "generic",
					NREntryPoint:    true,
					TransactionName: "WebTransaction/Go/GET /external",
				},
				UserAttributes: map[string]any{},
				AgentAttributes: map[string]any{
					"request.headers.host":   "localhost:8000",
					"code.filepath":          "go-agent/v3/examples/server/main.go",
					"request.uri":            "/external",
					"request.headers.accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
					"code.namespace":         "main",
					"request.method":         "GET",
					"code.lineno":            float64(160),
					"code.function":          "external",
				},
			},
		},
	}

	got, err := ParseSpanEventData([]byte(input))
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}
