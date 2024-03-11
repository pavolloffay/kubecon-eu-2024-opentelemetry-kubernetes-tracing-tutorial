# Manual instrumentation using the OpenTelemetry SDK

This tutorial section covers the manual instrumentation of a go application with the opentelemetry-sdk.

As a basis for the instrumentation we use [backend4](./app/backend4/main.go). To compile the application you need [go 1.22 or newer](https://go.dev/doc/install).

# Configure OpenTelemetry-go-sdk

```diff
func main() {
+	otelExporter, err := otlptracegrpc.New(context.Background())
+	if err != nil {
+		fmt.Printf("failed to create trace exporter: %s\n", err)
+		os.Exit(1)
+	}
+	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(otelExporter))
+	otel.SetTracerProvider(tp)
...
```


## Create and register a global trace provider

```diff
+var tracer = otel.GetTracerProvider().Tracer("github.com/kubecon-eu-2024/backend")
```

## Identifying critical path and operations for instrumentation

```diff
	mux.HandleFunc("GET /rolldice", func(w http.ResponseWriter, r *http.Request) {
+		var span trace.Span
+		ctx, span := tracer.Start(r.Context(), "rolldice")
+		defer span.End()
		player := "Anonymous player"
		if p := r.URL.Query().Get("player"); p != "" {
			player = p
		}
```

```diff
func causeError(ctx context.Context, rate int) error {
+	var span trace.Span
+	_, span = tracer.Start(ctx, "causeError")
+	defer span.End()

	randomNumber := rand.Intn(100)
+	span.AddEvent(fmt.Sprintf("random nr: %d", randomNumber))
	if randomNumber < rate {
		err := fmt.Errorf("internal server error")
+		span.RecordError(err)
		return err
	}
	return nil
}
```

```diff
func causeDelay(ctx context.Context, rate int) {
+	var span trace.Span
+	_, span = tracer.Start(ctx, "causeDelay")
+	defer span.End()
	randomNumber := rand.Intn(100)
+	span.AddEvent(fmt.Sprintf("random nr: %d", randomNumber))
	if randomNumber < rate {
		time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)
	}
}
```

## Configuring an OTLP exporter and setting the endpoint

```bash
docker run --rm -it -p 127.0.0.1:4317:4317 -p 127.0.0.1:16686:16686 -e COLLECTOR_OTLP_ENABLED=true -e LOG_LEVEL=debug  jaegertracing/all-in-one:latest
```

```bash
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 OTEL_SERVICE_NAME=go-backend go run app/backend4/main.go
```

## TODO: Publish container at ghcr.io/pavolloffay

Apply `backend2` drop-in replacement:
```bash
kubectl apply -f https://raw.githubusercontent.com/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial/main/backend/04-backend.yaml
```

Details

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend2-deployment
  namespace: tutorial-application
  labels:
    app: backend2
spec:
  template:
    metadata:
      labels:
        app: backend2
      annotations:
        prometheus.io/scrape: "true"
+        instrumentation.opentelemetry.io/inject-sdk: "true"
  template:
    spec:
      containers:
      - name: backend2
-        image: ghcr.io/pavolloffay/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial-backend2:latest
+        image: ghcr.io/frzifus/kubecon-eu-2024-opentelemetry-kubernetes-tracing-tutorial-backend4:latest
```

---

[Next steps](./05-sampling.md)
