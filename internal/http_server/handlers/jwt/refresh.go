package refresh

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "new-websocket-chat/internal/lib/api/response"
	jwtAuth "new-websocket-chat/internal/lib/jwt"
	"new-websocket-chat/internal/lib/logger/sl"
	"strconv"
)

type Response struct {
	resp.Response
	JWTAccessToken  string `json:"jwtAccessToken"`
	JWTRefreshToken string `json:"jwtRefreshToken"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.jwt.refresh.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		refreshToken := jwtAuth.ExtractToken(r)
		log.Info("refreshToken extracted", slog.String("refreshToken", refreshToken))

		claims, err := jwtAuth.ValidateToken(refreshToken)
		if err != nil {
			log.Error("Invalid refresh token", http.StatusUnauthorized)

			render.JSON(w, r, resp.Error("Invalid refresh token"))

			return
		}

		log.Info("refreshToken validated", slog.String("refreshToken", refreshToken))

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			log.Error("Failed to parse userID to int64", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to parse userID to int64"))

			return
		}
		log.Info("userID parsed from string to int64", slog.Int64("userID", userID))

		newAccessToken, newRefreshToken, err := jwtAuth.GenerateTokens(userID)
		if err != nil {
			log.Error("Failed to generate new access token", http.StatusInternalServerError)

			render.JSON(w, r, resp.Error("Failed to generate new access token"))

			return
		}
		log.Info("new access token generated", slog.String("newAccessToken", newAccessToken))

		responseOK(w, r, newAccessToken, newRefreshToken)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, generatedAccessToken string, generatedRefreshToken string) {
	render.JSON(w, r, Response{
		Response:        resp.OK(),
		JWTAccessToken:  generatedAccessToken,
		JWTRefreshToken: generatedRefreshToken,
	})
}
