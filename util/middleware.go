package util

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func JsonRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid content type", http.StatusBadRequest)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func LimitBody1Mb(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1000000)

		h.ServeHTTP(w, r)
	})
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			accessTokenCookie *http.Cookie
			err               error
		)
		accessTokenCookie, err = r.Cookie("accessToken")
		if err != nil {
			WriteApiMessage(w, "Missing access token", http.StatusUnauthorized, false)
			return
		}

		var (
			jwtClaims  jwt.MapClaims
			statusCode int
		)
		jwtClaims, statusCode, err = VerifyJwtToken(accessTokenCookie.Value)
		if err != nil {
			WriteApiMessage(w, err.Error(), statusCode, false)
			return
		}

		// Would have to typecast these values again, whereever they're going to be used anyway
		// Directly going to call ParseJwtClaims in those locations
		// IGNORE
		var userDetails map[string]interface{}
		userDetails, statusCode, err = ParseJwtClaims(jwtClaims)
		if err != nil {
			WriteApiMessage(w, err.Error(), statusCode, false)
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userDetails", userDetails)))
	})
}
