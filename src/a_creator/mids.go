package main

import (
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

		next.ServeHTTP(w, r)
	})
}

func fakeAuthMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract from headers Bearer {token}
		// validate {token}
		// PROFIT !!!!

		next.ServeHTTP(w, r)
	})
}
