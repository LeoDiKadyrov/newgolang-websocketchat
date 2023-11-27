package refresh

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "new-websocket-chat/internal/lib/api/response"
	"new-websocket-chat/internal/lib/logger/sl"
	"strconv"
)

type Response struct {
	resp.Response
	JWTAccessToken  string `json:"jwtAccessToken"`
	JWTRefreshToken string `json:"jwtRefreshToken"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.37.1 --name=TokenService
type TokenService interface {
	ExtractToken(r *http.Request) (string, error)
	ValidateToken(tokenString string) (*jwt.StandardClaims, error)
	GenerateTokens(userID int64) (accessTokenString string, refreshTokenString string, err error)
}

func New(log *slog.Logger, tokenService TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.jwt.refresh.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		refreshToken, err := tokenService.ExtractToken(r)
		if err != nil {
			log.Error("Failed to extract token from request", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to extract token"))

			return
		}

		log.Info("refreshToken extracted", slog.String("refreshToken", refreshToken))

		claims, err := tokenService.ValidateToken(refreshToken)
		if err != nil {
			log.Error("Invalid refresh token", sl.Err(err))

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

		newAccessToken, newRefreshToken, err := tokenService.GenerateTokens(userID)
		if err != nil {
			log.Error("Failed to generate new access token", sl.Err(err))

			render.JSON(w, r, resp.Error("Failed to generate new access token"))

			return
		}
		log.Info("new access token generated", slog.String("newAccessToken", newAccessToken)) // TODO: remove due to security issues
		log.Info("new refresh token generated", slog.String("newAccessToken", newRefreshToken))

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
