package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.GetTracerProvider().Tracer("github.com/kubecon-eu-2024/backend")

var (
	rollCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dice_roll_count",
			Help: "How often the dice was rolled",
		},
	)

	numbersCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dice_numbers_count",
			Help: "How often each number of the dice was rolled",
		},
		[]string{"number"},
	)
)

func init() {
	prometheus.MustRegister(rollCounter)
	prometheus.MustRegister(numbersCounter)
}

func main() {
	otelExporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		fmt.Printf("failed to create trace exporter: %s\n", err)
		os.Exit(1)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(otelExporter))
	otel.SetTracerProvider(tp)

	v, ok := os.LookupEnv("ERROR_RATE")
	if !ok {
		v = "0"
	}
	rateError, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	v, ok = os.LookupEnv("RATE_HIGH_DELAY")
	if !ok {
		v = "0"
	}
	rateDelay, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rolldice", func(w http.ResponseWriter, r *http.Request) {
		var span trace.Span
		ctx, span := tracer.Start(r.Context(), "rolldice")
		defer span.End()
		player := "Anonymous player"
		if p := r.URL.Query().Get("player"); p != "" {
			player = p
		}
		max := 8
		if fmt.Sprintf("%x", sha256.Sum256([]byte(player))) == "f4b7c19317c929d2a34297d6229defe5262fa556ef654b600fc98f02c6d87fdc" {
			max = 8
		} else {
			max = 6
		}
		result := doRoll(ctx, max)
		causeDelay(ctx, rateDelay)
		if err := causeError(ctx, rateError); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resStr := strconv.Itoa(result)
		rollCounter.Inc()
		numbersCounter.WithLabelValues(resStr).Inc()
		if _, err := w.Write([]byte(resStr)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

	})
	mux.HandleFunc("GET /metrics", promhttp.Handler().ServeHTTP)
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func causeError(ctx context.Context, rate int) error {
	var span trace.Span
	_, span = tracer.Start(ctx, "causeError")
	defer span.End()

	randomNumber := rand.Intn(100)
	span.AddEvent(fmt.Sprintf("random nr: %d", randomNumber))
	if randomNumber < rate {
		err := fmt.Errorf("internal server error")
		span.RecordError(err)
		return err
	}
	return nil
}

func causeDelay(ctx context.Context, rate int) {
	var span trace.Span
	_, span = tracer.Start(ctx, "causeDelay")
	defer span.End()
	randomNumber := rand.Intn(100)
	span.AddEvent(fmt.Sprintf("random nr: %d", randomNumber))
	if randomNumber < rate {
		time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)
	}
}

func doRoll(_ context.Context, max int) int {
	return rand.Intn(max) + 1
}
