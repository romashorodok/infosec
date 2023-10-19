package httputils

import (
	"fmt"
	"net/http"

	"github.com/romashorodok/infosec/pkg/auth"
	"log/slog"
)

func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.WithTokenPayload(r.Context())
			if err != nil {
				logger.Error("logger unable get user identity", "err", err)
			}
			logger.Info(fmt.Sprintf("request on %+v", r.URL), "user", fmt.Sprintf("%+v", token))

			next.ServeHTTP(w, r)
		})
	}
}
