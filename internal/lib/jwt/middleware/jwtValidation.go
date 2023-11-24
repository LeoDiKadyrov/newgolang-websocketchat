package jwtAuth

import (
	"context"
	"net/http"
	jwtAuth "new-websocket-chat/internal/lib/jwt"
)

func TokenAuthMiddleware(next http.Handler) http.Handler {
	const op = "lib.jwt.middleware.TokenAuthMiddleware"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := jwtAuth.ExtractToken(r)

		claims, err := jwtAuth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
