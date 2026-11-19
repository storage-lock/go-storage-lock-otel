# Storage Lock OpenTelemetry

把 storage-lock 事件记录为 OpenTelemetry span。监听器只观测，不改变加锁结果。

```go
listener := otel.NewTracer(otel.Tracer("storage-lock"))
options.AddEventListeners(listener)
```

每个事件创建一个 span，属性为 `lock.id`、`lock.owner_id`、`lock.storage`，每个 action 记录为一个 span event。

## 与 Prometheus 指标的区别

这里**故意**把 `lock.id` 和 `lock.owner_id` 作为 span 属性——追踪是按单次调用存储的，
本来就用来定位"这一次请求发生了什么"。这和 `go-storage-lock-prometheus`
的低基数原则不冲突，两者解决不同问题：

| | Prometheus 指标 | OTel span |
|---|---|---|
| 存储方式 | 聚合时序 | 单条记录 |
| 含 lock_id | ❌ 会打爆时序库 | ✅ 追踪需要 |
| 回答 | 整体趋势 | 这一次调用 |

不要用 span 属性做聚合统计，也不要把 lock_id 加进 Prometheus label。
