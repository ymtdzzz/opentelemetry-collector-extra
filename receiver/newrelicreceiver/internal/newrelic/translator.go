package newrelic

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/collector/semconv/v1.27.0"
)

func ToTraces(conn *Connection, nrSpans *SpanEvent) (ptrace.Traces, error) {
	results := ptrace.NewTraces()
	rs := results.ResourceSpans().AppendEmpty()
	rs.SetSchemaUrl(semconv.SchemaURL)
	for k, v := range map[string]string{
		semconv.AttributeProcessPID:           strconv.Itoa(conn.Pid),
		semconv.AttributeHostName:             conn.Host,
		semconv.AttributeTelemetrySDKLanguage: conn.Language,
		semconv.AttributeTelemetrySDKVersion:  conn.Version,
		semconv.AttributeTelemetrySDKName:     "NewRelic",
		semconv.AttributeServiceName:          conn.AppName[0],
	} {
		if v != "" {
			rs.Resource().Attributes().PutStr(k, v)
		}
	}

	in := rs.ScopeSpans().AppendEmpty()
	in.Scope().SetName("NewRelic")
	in.Scope().SetVersion(conn.Version)

	for _, span := range nrSpans.Spans {
		newSpan := in.Spans().AppendEmpty()

		traceID, err := convertTraceID(span.TraceID)
		if err != nil {
			return results, err
		}
		newSpan.SetTraceID(traceID)
		spanID, err := convertSpanID(span.SpanID)
		if err != nil {
			return results, err
		}
		newSpan.SetSpanID(spanID)
		start := time.UnixMilli(int64(span.Timestamp))
		dur := time.Duration(span.Duration * float64(time.Second))
		newSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
		newSpan.SetEndTimestamp(pcommon.NewTimestampFromTime(start.Add(dur)))
		if span.ParentID != "" {
			parentID, err := convertSpanID(span.ParentID)
			if err != nil {
				return results, err
			}
			newSpan.SetParentSpanID(parentID)
		}
		newSpan.SetName(span.Name)

		// TODO: handle other kind, component and category?
		switch span.SpanKind {
		case "server":
			newSpan.SetKind(ptrace.SpanKindServer)
		case "client":
			newSpan.SetKind(ptrace.SpanKindClient)
		default:
			newSpan.SetKind(ptrace.SpanKindUnspecified)
		}

		// TODO: convert to OTel semantic attributes
		// see: https://github.com/newrelic/go-agent/blob/master/v3/newrelic/attributes.go
		setAttributes(&newSpan, span.UserAttributes)
		setAttributes(&newSpan, span.AgentAttributes)
	}

	return results, nil
}

func setAttributes(span *ptrace.Span, attrs map[string]any) {
	for k, anyV := range attrs {
		switch v := anyV.(type) {
		// NOTE: all number values are parsed to float64
		case string:
			span.Attributes().PutStr(k, v)
		case bool:
			span.Attributes().PutBool(k, v)
		case float64:
			span.Attributes().PutDouble(k, v)
		default:
			span.Attributes().PutStr(k, fmt.Sprintf("%v", v))
		}
	}
}

func convertTraceID(nrTraceID string) (pcommon.TraceID, error) {
	var traceID [16]byte

	if len(nrTraceID) != 32 {
		return traceID, errors.New("invalid Trace ID: must be exactly 32 hexdecimal characters")
	}

	bytes, err := hex.DecodeString(nrTraceID)
	if err != nil {
		return traceID, fmt.Errorf("failed to decode Trace ID: %v", err)
	}

	copy(traceID[:], bytes)

	return traceID, nil
}

func convertSpanID(nrSpanID string) (pcommon.SpanID, error) {
	var spanID [8]byte

	if len(nrSpanID) != 16 {
		return spanID, errors.New("invalid Span ID: must be exactly 16 hexdecimal characters")
	}

	bytes, err := hex.DecodeString(nrSpanID)
	if err != nil {
		return spanID, fmt.Errorf("failed to decode Trace ID: %v", err)
	}

	copy(spanID[:], bytes)

	return spanID, nil
}
