package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// load ENV vars
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// initialize tracer
	tracer, closer, err := NewJaegerOpentracingTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// build dependency graph
	repo := NewRepo()
	publisher, err := NewPublisher()
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.amqpChan.Close()
	defer publisher.conn.Close()

	http.Handle("/order/create", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := new(CreateOrderRequest)
		if err := json.Unmarshal(body, req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		order, err := CreateOrder(r.Context(), req, repo, publisher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}))

	if err := http.ListenAndServe(os.Getenv("APP_ADDR"), entryPointMid(fakeAuthMid(http.DefaultServeMux))); err != nil {
		log.Fatal(err)
	}
}
