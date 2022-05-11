package main

import (
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io"
	"log"
	"net/http"
)

func main() {
	// initialize tracer
	tracer, closer, err := NewJaegerOpentracingTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// define routes
	http.HandleFunc("/process_bill", func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "controller: process_bill")
		defer span.Finish()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			ext.LogError(span, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := new(ProcessBillRequest)
		if err := json.Unmarshal(body, req); err != nil {
			ext.LogError(span, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ProcessBill(ctx, req); err != nil {
			ext.LogError(span, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(":4001", entryPointMid(xAPiMid(http.DefaultServeMux))); err != nil {
		log.Fatal(err)
	}
}
