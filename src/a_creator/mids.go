package main

import (
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

		span, ctx := opentracing.StartSpanFromContext(r.Context(), r.URL.Path)
		defer span.Finish()

		// extract ip address
		remoteAddr := r.Header.Get("X-Real-IP")
		if remoteAddr == "" {
			remoteAddr = r.Header.Get("X-Forwarded-For")
		}
		if remoteAddr == "" {
			remoteAddr = r.RemoteAddr
		}

		ext.PeerAddress.Set(span, remoteAddr)
		ext.HTTPUrl.Set(span, r.URL.Path)
		ext.HTTPMethod.Set(span, r.Method)

		rid := uuid.NewString()
		span.SetBaggageItem("request_id", rid)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func fakeAuthMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "fakeAuthMid")
		defer span.Finish()

		// extract from headers Bearer {token}
		// validate somehow {token}
		// if success let it pass through chain, otherwise log error and return

		fakeSID := uuid.NewString()
		span.SetBaggageItem("subject_id", fakeSID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
