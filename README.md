# Storage Lock OpenTelemetry

把 storage-lock 事件记录为 OpenTelemetry span。监听器只观测，不改变加锁结果。

```go
listener := otel.NewTracer(otel.Tracer("storage-lock"))
options.AddEventListeners(listener)
```

每个事件创建一个 span，属性为 `lock.id`、`lock.owner_id`、`lock.storage`，每个 action 记录为一个 span event。
