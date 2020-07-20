# promezeus
Abstraction around prometheus module

## prometrics

```go
import (
    "github.com/eduardbcom/promezeus/prometrics"
)

metricsServer, err := prometrics.New()
if err != nil {
    return nil, err
}

// IMPORTANT:
// 9100 port is hardcoded.
metricsServer.Listen()

// ...

prometrics.RegisterGauge(
    // metric key used in order to identify unique gauge (can be blank though)
    // (metricKey, metricName) must be a uniq pair in order to avoid collisions. 
    "metricKey", 
    prometrics.GaugeType,
    // metric name is a string you will see within /metrics report
    "metricName",
    "Too long description for metricName",
    // labels (can be nil)
    map[string]string{"id": "1"},
    func(labels prometrics.Labels) float64 { /* value gette */ return 1.0 }
)

# HELP metricName Too long description for metricName
# TYPE metricName gauge
metricName{id="1"} 1.0

// ...

metricsServer.StopListen()
```

## promquery

```go
import (
    "github.com/eduardbcom/promezeus/promquery"
)

promQuery, err := promquery.New(
    &promquery.Config{
        PromAPI: "http://your-prometheus-url.com",
    },
)
if err != nil {
    return nil, err
}

ctx := context.TODO()

promQuery.Query(ctx, "metricName{id='1'}")
```