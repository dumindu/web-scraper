package middleware

import (
	"net/http"
	"strings"

	"web-scraper.dev/internal/utils/ctxutil"
	"web-scraper.dev/internal/utils/jwtutil"
)

func JwtAuthentication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) < 7 || strings.ToUpper(tokenString[0:6]) != "BEARER" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "no token found"}`))
			return
		}

		tokenString = tokenString[7:]

		claims, err := jwtutil.ClaimsFromAccessToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid token"}`))
			return
		}

		ctxUser := claims.ToCtxUser()
		if ctxUser.Email == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid token"}`))
			return
		}
		ctx := ctxutil.SetUser(r.Context(), ctxUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
