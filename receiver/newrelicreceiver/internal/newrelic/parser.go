package newrelic

import (
	"encoding/json"
	"errors"
	"fmt"
)

type mixedData []interface{}

func ParseConnectionData(payload []byte) (*Connection, error) {
	var conns []Connection

	if err := json.Unmarshal(payload, &conns); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}
	if len(conns) == 0 {
		return nil, errors.New("parsed connections has no elements")
	}

	return &conns[0], nil
}

func ParseSpanEventData(payload []byte) (*SpanEvent, error) {
	var mixed mixedData

	if err := json.Unmarshal(payload, &mixed); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	if len(mixed) != 3 {
		return nil, fmt.Errorf("parsed payload doesn't have 3 elements. it has %d elements", len(mixed))
	}

	runID, ok := mixed[0].(string)
	if !ok {
		return nil, errors.New("the first element of parsed payload is not string")
	}

	einfoJSON, _ := json.Marshal(mixed[1])
	var einfo EventInfo
	err := json.Unmarshal(einfoJSON, &einfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the second element: %w", err)
	}

	spans, ok := mixed[2].([]interface{})
	if !ok {
		return nil, errors.New("the third element of parsed payload is not array")
	}

	var ss []*Span
	for _, span := range spans {
		s, ok := span.([]interface{})
		if !ok || len(s) != 3 {
			return nil, fmt.Errorf("parsed span data doesn't have 3 elements. actual element count: %d", len(s))
		}
		spanJSON, _ := json.Marshal(s[0])
		var span ParsedSpan
		err := json.Unmarshal(spanJSON, &span)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the first element: %w", err)
		}

		userAttrs, ok := s[1].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to parse the second element as map[string]any: %w", err)
		}

		agentAttrs, ok := s[2].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to parse the third element as map[string]any: %w", err)
		}

		ss = append(ss, &Span{
			ParsedSpan:      &span,
			UserAttributes:  userAttrs,
			AgentAttributes: agentAttrs,
		})
	}

	return &SpanEvent{
		RunID:     runID,
		EventInfo: &einfo,
		Spans:     ss,
	}, nil
}

func ParseTransactionEventData(payload []byte) (*TransactionEvent, error) {
	var mixed mixedData

	if err := json.Unmarshal(payload, &mixed); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	if len(mixed) != 3 {
		return nil, fmt.Errorf("parsed payload doesn't have 3 elements. it has %d elements", len(mixed))
	}

	runID, ok := mixed[0].(string)
	if !ok {
		return nil, errors.New("the first element of parsed payload is not string")
	}

	einfoJSON, _ := json.Marshal(mixed[1])
	var einfo EventInfo
	err := json.Unmarshal(einfoJSON, &einfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the second element: %w", err)
	}

	txs, ok := mixed[2].([]interface{})
	if !ok {
		return nil, errors.New("the third element of parsed payload is not array")
	}

	var ts []*Transaction
	for _, tx := range txs {
		s, ok := tx.([]interface{})
		if !ok || len(s) != 3 {
			return nil, fmt.Errorf("parsed span data doesn't have 3 elements. actual element count: %d", len(s))
		}
		spanJSON, _ := json.Marshal(s[0])
		var t ParsedTransaction
		err := json.Unmarshal(spanJSON, &t)
		if err != nil {
			return nil, fmt.Errorf("failed to parse the first element: %w", err)
		}

		userAttrs, ok := s[1].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to parse the second element as map[string]any: %w", err)
		}

		agentAttrs, ok := s[2].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to parse the third element as map[string]any: %w", err)
		}

		ts = append(ts, &Transaction{
			ParsedTransaction: &t,
			UserAttributes:    userAttrs,
			AgentAttributes:   agentAttrs,
		})
	}

	return &TransactionEvent{
		RunID:        runID,
		EventInfo:    &einfo,
		Transactions: ts,
	}, nil
}
