package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

func entryPointMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create the span referring to the RPC client if available.
		// If spanContext == nil, a root span will be created.
		span := opentracing.StartSpan("entryPointMid", ext.RPCServerOption(spanContext))
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(r.Context(), span)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func xAPiMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "xAPiMid")
		defer span.Finish()

		key := r.Header.Get("x-api-key")
		if key == "" {
			errMsg := "'x-api-key' is empty"
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		// only for demo purpose
		// NEVER log `secret`s
		keyPart := fmt.Sprintf("%s***", key[:len(key)/2])
		span.LogFields(log.String("key", keyPart))

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
