# promezeus
Abstraction around prometheus module

# prometrics

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

// example with plain gauge

prometrics.Register(
    "metricName",
    []{"id"}
     "Too long description for metricName",
)
prometrics.Inc("metricName", map[string]string{"id": "1"}

# HELP metricName Too long description for metricName
# TYPE metricName gauge
metricName{id="1"} 1


// example with collector

prometrics.RegisterWithCollector(
    "metricName",
    map[string]string{"id": "1"},
    "Too long description for metricName",
    // labels (can be nil)
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

# TODO:
- add support for counter
- unit/func tests