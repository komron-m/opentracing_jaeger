package main

import (
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

func entryPointMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.Write(nil)
			return
		}

		span, ctx := opentracing.StartSpanFromContext(r.Context(), "check-cors(entrypoint)")
		defer span.Finish()

		// tags
		span.SetTag("is_mobile", true)

		span.LogKV("http.method", r.Method)
		span.LogKV("ip.addr", r.RemoteAddr)

		requestId := uuid.NewString()
		span.SetBaggageItem("request_id", requestId)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func fakeAuthMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract from headers Bearer {token}
		// validate {token}
		// PROFIT !!!!

		span, ctx := opentracing.StartSpanFromContext(r.Context(), "fakeAuthMid")
		defer span.Finish()

		userID := uuid.NewString()
		span.SetBaggageItem("user_id", userID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
