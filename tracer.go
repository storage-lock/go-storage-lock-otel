// Package otel 把 storage-lock 的事件流记录为 OpenTelemetry span。
package otel

import (
	"context"

	"github.com/storage-lock/go-events"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Tracer 是可选的事件监听器，不参与加锁决策。
type Tracer struct {
	tracer oteltrace.Tracer
}

var _ events.Listener = &Tracer{}

const TracerName = "storage-lock-otel"

// NewTracer 使用调用方 tracer provider 中的 tracer。provider 的生命周期由调用方管理。
func NewTracer(tracer oteltrace.Tracer) *Tracer {
	return &Tracer{tracer: tracer}
}

func (x *Tracer) Name() string { return TracerName }

// On 为每个事件创建一个 span，并把 action 名称和错误记入属性。
func (x *Tracer) On(ctx context.Context, e *events.Event) {
	if x == nil || x.tracer == nil || e == nil {
		return
	}
	names := []string{"storage-lock", "create", "lock", "unlock", "watchdog"}
	name := "storage-lock"
	if int(e.EventType) >= 0 && int(e.EventType) < len(names) {
		name = names[e.EventType]
	}
	_, span := x.tracer.Start(ctx, name, oteltrace.WithAttributes(
		attribute.String("lock.id", e.LockId),
		attribute.String("lock.owner_id", e.OwnerId),
		attribute.String("lock.storage", e.StorageName),
	))
	defer span.End()
	for _, action := range e.Actions {
		if action == nil {
			continue
		}
		span.AddEvent(action.Name)
		if action.Err != nil {
			span.RecordError(action.Err)
			span.SetStatus(codes.Error, action.Err.Error())
		}
	}
	if e.Err != nil {
		span.RecordError(e.Err)
		span.SetStatus(codes.Error, e.Err.Error())
	}
}
