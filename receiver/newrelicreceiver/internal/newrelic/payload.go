package newrelic

// Connection is information of current connection. This is used to populate Resource and Instrumentation scope attributes
// see: https://github.com/newrelic/go-agent/blob/82c8f8440ca84eb68e08248d877fa1d0b55da333/v3/newrelic/config.go#L935-L951
type Connection struct {
	Pid      int    `json:"pid"`
	Language string `json:"language"`
	Version  string `json:"agent_version"`
	Host     string `json:"host"`
	// HostDisplayName  string   `json:"display_host,omitempty"`
	// Settings         any      `json:"settings"`
	AppName []string `json:"app_name"`
	// HighSecurity     bool     `json:"high_security"`
	// Labels           any      `json:"labels,omitempty"`
	// Environment      any      `json:"environment"`
	// Identifier       string   `json:"identifier"`
	// Util             any      `json:"utilization"`
	// SecurityPolicies any      `json:"security_policies,omitempty"`
	// Metadata         any      `json:"metadata"`
	// EventData        any      `json:"event_harvest_config"`
}

type EventInfo struct {
	ReservoirSize uint64 `json:"reservoir_size"`
	EventsSeen    uint64 `json:"events_seen"`
}

// ParsedSpan is span from NewRelic agent
// see: https://github.com/newrelic/go-agent/blob/82c8f8440ca84eb68e08248d877fa1d0b55da333/v3/newrelic/span_events.go#L22-L42
type ParsedSpan struct {
	Type            string  `json:"type"`
	TraceID         string  `json:"traceId"`
	SpanID          string  `json:"guid"`
	ParentID        string  `json:"parentId,omitempty"`
	TransactionID   string  `json:"transactionId"`
	Sampled         bool    `json:"sampled"`
	Priority        float64 `json:"priority"`
	Timestamp       uint64  `json:"timestamp"` // millis
	Duration        float64 `json:"duration"`  // seconds
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	NREntryPoint    bool    `json:"nr.entryPoint,omitempty"`
	Component       string  `json:"component,omitempty"`
	SpanKind        string  `json:"span.kind,omitempty"`
	TrustedParentID string  `json:"trustedParentId,omitempty"`
	TracingVendors  string  `json:"tracingVendors,omitempty"`
	TransactionName string  `json:"transaction.name,omitempty"`
}

type Span struct {
	*ParsedSpan
	UserAttributes  map[string]any
	AgentAttributes map[string]any
}

// SpanEvent is full span event from NewRelic agent
type SpanEvent struct {
	RunID string
	*EventInfo
	Spans []*Span
}

// ParsedTransaction is transaction from NewRelic agent
// see: https://github.com/newrelic/go-agent/blob/82c8f8440ca84eb68e08248d877fa1d0b55da333/v3/newrelic/txn_events.go#L15
type ParsedTransaction struct {
	Type              string  `json:"type"`
	Name              string  `json:"name"`
	Timestamp         uint64  `json:"timestamp"`
	NRApdexPerfZone   string  `json:"nr.apdexPerfZone,omitempty"`
	HasError          bool    `json:"error"`
	Duration          float64 `json:"duration"`                // seconds
	QueueDuration     float64 `json:"queueDuration,omitempty"` // seconds
	ExternalCallCount int64   `json:"externalCallCount,omitempty"`
	ExternalDuration  float64 `json:"externalDuration,omitempty"` // seconds
	DatabaseCallCount int64   `json:"databaseCallCount,omitempty"`
	DatabaseDuration  float64 `json:"databaseDuration,omitempty"` // seconds
	// TODO: Synthetics related information
	TotalTime float64 `json:"totalTime"`
	// TODO: parent information
	TransactionID string  `json:"guid,omitempty"`
	TraceID       string  `json:"traceId,omitempty"`
	Priority      float64 `json:"priority,omitempty"`
	Sampled       bool    `json:"sampled,omitempty"`
}

type Transaction struct {
	*ParsedTransaction
	UserAttributes  map[string]any
	AgentAttributes map[string]any
}

// TransactionEvent is full transaction event from NewRelic agent
type TransactionEvent struct {
	RunID string
	*EventInfo
	Transactions []*Transaction
}
