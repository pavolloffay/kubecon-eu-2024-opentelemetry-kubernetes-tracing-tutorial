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
)

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
		result := doRoll(r.Context(), max)
		if err := causeError(r.Context(), rateError); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		causeDelay(r.Context(), rateDelay)
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

func causeError(_ context.Context, rate int) error {
	randomNumber := rand.Intn(100)
	if randomNumber < rate {
		return fmt.Errorf("internal server error")
	}
	return nil
}

func causeDelay(_ context.Context, rate int) {
	randomNumber := rand.Intn(100)
	if randomNumber < rate {
		time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)
	}
}

func doRoll(_ context.Context, max int) int {
	return rand.Intn(max) + 1
}
