package newrelicreceiver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"net/http"

	"github.com/ymtdzzz/opentelemetry-collector-extra/receiver/newrelicreceiver/internal/newrelic"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"google.golang.org/grpc"
)

type newrelicReceiver struct {
	address string
	config  *Config
	params  receiver.Settings
	conns   map[string]*newrelic.Connection

	nextTracesConsumer consumer.Traces

	server     *grpc.Server
	httpServer *http.Server
	tReceiver  *receiverhelper.ObsReport
}

type PreconnectReplyResponse struct {
	Preconnect PreconnectReply `json:"return_value"`
}

// PreconnectReply is response from NewRelic's preconnect endpoint.
// see: https://github.com/newrelic/go-agent/blob/82c8f8440ca84eb68e08248d877fa1d0b55da333/v3/internal/connect_reply.go#L26-L30
type PreconnectReply struct {
	Collector        string           `json:"redirect_host"`
	SecurityPolicies SecurityPolicies `json:"security_policies"`
}

// SecurityPolicies contains the security policies.
// see: https://github.com/newrelic/go-agent/blob/82c8f8440ca84eb68e08248d877fa1d0b55da333/v3/internal/security_policies.go#L16-L22
type SecurityPolicies struct {
	RecordSQL                 securityPolicy `json:"record_sql"`
	AttributesInclude         securityPolicy `json:"attributes_include"`
	AllowRawExceptionMessages securityPolicy `json:"allow_raw_exception_messages"`
	CustomEvents              securityPolicy `json:"custom_events"`
	CustomParameters          securityPolicy `json:"custom_parameters"`
}

type securityPolicy struct {
	EnabledVal *bool `json:"enabled"`
}

type ConnectReplyResponse struct {
	Connect ConnectReply `json:"return_value"`
}

// ConnectReply is response from NewRelic's connect endpoint.
// see: https://github.com/newrelic/go-agent/blob/82c8f8440ca84eb68e08248d877fa1d0b55da333/v3/internal/connect_reply.go#L34
type ConnectReply struct {
	RunID string `json:"agent_run_id"`
}

func newNewRelicReceiver(config *Config, params receiver.Settings) (component.Component, error) {
	instance, err := receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{
		LongLivedCtx:           false,
		ReceiverID:             params.ID,
		Transport:              "http",
		ReceiverCreateSettings: params,
	})
	if err != nil {
		return nil, err
	}

	return &newrelicReceiver{
		config:    config,
		params:    params,
		conns:     map[string]*newrelic.Connection{},
		tReceiver: instance,
	}, nil
}

func (nrr *newrelicReceiver) Start(ctx context.Context, host component.Host) error {
	nrmux := http.NewServeMux()
	nrmux.HandleFunc("/agent_listener/invoke_raw_method", nrr.handleInvokeRawMethod)

	scfg := confighttp.NewDefaultServerConfig()
	scfg.Endpoint = "127.0.0.1:8080"
	scfg.TLSSetting.CertFile = nrr.config.ServerCert
	scfg.TLSSetting.KeyFile = nrr.config.ServerKey
	var err error
	nrr.httpServer, err = scfg.ToServer(
		ctx,
		host,
		nrr.params.TelemetrySettings,
		nrmux,
	)
	if err != nil {
		return fmt.Errorf("failed to create server definition: %w", err)
	}
	hln, err := scfg.ToListener(ctx)
	if err != nil {
		return fmt.Errorf("failed to create datadog listener: %w", err)
	}

	nrr.address = hln.Addr().String()

	go func() {
		if err := nrr.httpServer.Serve(hln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(fmt.Errorf("error starting datadog receiver: %w", err)))
		}
	}()

	return nil
}

func (nrr *newrelicReceiver) Shutdown(ctx context.Context) error {
	nrr.server.GracefulStop()
	return nil
}

func (nrr *newrelicReceiver) handleInvokeRawMethod(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	if q == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	method := q.Get("method")

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	if method == "preconnect" {
		resp := &PreconnectReplyResponse{
			PreconnectReply{
				Collector: "localhost:8080",
			},
		}
		b, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s := `{"status":"500 Internal Server Error"}`
			if _, err := w.Write([]byte(s)); err != nil {
				log.Println(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			log.Println(err)
		}
	} else if method == "connect" {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("error reading body: %v\n", err)
			return
		}
		conn, err := newrelic.ParseConnectionData(bodyBytes)
		if err != nil {
			log.Printf("error parsing payload: %v\n", err)
			return
		}

		aname := conn.AppName[0]
		nrr.conns[aname] = conn

		resp := &ConnectReplyResponse{
			ConnectReply{
				RunID: aname,
			},
		}
		b, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s := `{"status":"500 Internal Server Error"}`
			if _, err := w.Write([]byte(s)); err != nil {
				log.Println(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			log.Println(err)
		}
	} else if method == "span_event_data" {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("error reading body: %v\n", err)
			return
		}
		span, err := newrelic.ParseSpanEventData(bodyBytes)
		if err != nil {
			log.Printf("error parsing payload: %v\n", err)
			return
		}
		conn, ok := nrr.conns[span.RunID]
		if !ok {
			log.Println("error connection not found")
			return
		}
		otelspan, err := newrelic.ToTraces(conn, span)
		if err != nil {
			log.Printf("error convert NewRelic spans to OTel traces: %v\n", err)
			return
		}
		// TODO: error handling
		_ = nrr.nextTracesConsumer.ConsumeTraces(req.Context(), otelspan)

		w.WriteHeader(http.StatusOK)
	} else if method == "analytic_event_data" {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("error reading body: %v\n", err)
			return
		}
		tx, err := newrelic.ParseTransactionEventData(bodyBytes)
		if err != nil {
			log.Printf("error parsing payload: %v\n", err)
			return
		}
		log.Printf("Request payload (analytic_event_data):\n%#v\n", tx)
	}
}
