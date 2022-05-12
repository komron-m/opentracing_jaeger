package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

	// define http.routes
	http.Handle("/home", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, _ := opentracing.StartSpanFromContext(r.Context(), "greeting")
		defer span.Finish()

		w.Write([]byte("hello opentracing world!"))
	}))

	http.Handle("/order/create", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "order-create")
		defer span.Finish()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			ext.LogError(span, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := new(CreateOrderRequest)
		err = json.Unmarshal(body, req)
		if err != nil {
			ext.LogError(span, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		order, err := CreateOrder(ctx, req, repo, publisher)
		if err != nil {
			ext.LogError(span, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}))

	http.Handle("/order/list", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, _ := opentracing.StartSpanFromContext(r.Context(), "order-list")
		defer span.Finish()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repo.orders)
	}))

	if err := http.ListenAndServe(os.Getenv("APP_ADDR"), entryPointMid(fakeAuthMid(http.DefaultServeMux))); err != nil {
		log.Fatal(err)
	}
}
