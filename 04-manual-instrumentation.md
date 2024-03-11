# Manual instrumentation using the OpenTelemetry SDK

# Register Tracer

```go
var tracer = otel.GetTracerProvider().Tracer("github.com/kubecon-eu-2024/backend")
```


# init

```go
	var otlpAddr = flag.String("otlp-grpc", "", "default otlp/gRPC address, by default disabled. Example value: localhost:4317")
	flag.Parse()
	if *otlpAddr != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		grpcOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock()}
		conn, err := grpc.DialContext(ctx, *otlpAddr, grpcOptions...)
		if err != nil {
			fmt.Printf("failed to create gRPC connection to collector: %s\n", err)
			os.Exit(1)
		}
		defer conn.Close()

		// Set up a trace exporter
		otelExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			fmt.Printf("failed to create trace exporter: %s\n", err)
			os.Exit(1)
		}
		tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(otelExporter))
		otel.SetTracerProvider(tp)
	}
```

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

---

[Next steps](./05-sampling.md)
