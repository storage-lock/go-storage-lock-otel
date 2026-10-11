package otel

import (
	"context"
	"errors"
	"testing"

	"github.com/storage-lock/go-events"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTracerRecordsLockSpan(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	listener := NewTracer(provider.Tracer("test"))
	if listener.Name() != TracerName {
		t.Fatal(listener.Name())
	}

	event := events.NewEvent("orders").SetOwnerId("worker-1").SetStorageName("redis").SetType(events.EventTypeLock)
	event.AddAction(events.NewAction("StorageLock.Lock.Begin.Success"))
	event.AddAction(events.NewAction("StorageLock.Lock.Begin.Error").SetErr(errors.New("busy")))
	listener.On(context.Background(), event)
	listener.On(context.Background(), nil)

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("spans=%d", len(spans))
	}
	if spans[0].Name() != "lock" || len(spans[0].Events()) < 2 {
		t.Fatalf("span=%s events=%d", spans[0].Name(), len(spans[0].Events()))
	}
}
