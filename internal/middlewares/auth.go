package middlewares

import (
	"context"
	"net/http"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/utils"
)

type contextKey string

const contextKeyUserID contextKey = "userID"

// RequireAuth validates "Bearer <token>" and stores the userID in context.
func RequireAuth(jwtSecret string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := utils.GetBearerToken(r.Header)
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			userID, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				utils.WriteJSONError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the authenticated userID from context.
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(contextKeyUserID).(int64)
	return userID, ok
}
