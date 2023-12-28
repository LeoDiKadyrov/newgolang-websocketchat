package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	resp "new-websocket-chat/internal/lib/api/response"
	"new-websocket-chat/internal/lib/encryption"
	jwtAuth "new-websocket-chat/internal/lib/jwt"
	"new-websocket-chat/internal/lib/logger/sl"
	"new-websocket-chat/internal/storage"
)

// Request defines the required information to create a new user.
type Request struct {
	Username string `json:"username" validate:"required,min=4,max=24"` // Username of the user
	Email    string `json:"email" validate:"required,email"` // Email of the user
	Password string `json:"password" validate:"required,min=8,max=24,containsany=!@#?"` // Password of the user
}

// Response defines the response payload for the user creation request.
type Response struct {
	resp.Response // Embedding the common response struct
	Username        string `json:"username,omitempty"` // Username that was registered
	JWTAccessToken  string `json:"jwtAccessToken"` // Access JWT token for the user
	JWTRefreshToken string `json:"jwtRefreshToken"` // Refresh JWT token for the user
}

//go:generate go run github.com/vektra/mockery/v2@v2.37.1 --name=UserSaver
type UserSaver interface {
	SaveUser(username string, email string, password string) (int64, error)
}

// @Summary Create user
// @Description Create a new user in the system.
// @Tags user
// @Accept json
// @Produce json
// @Param request body Request true "User Registration Data"
// @Success 200 {object} Response "Successfully registered user and generated JWT tokens"
// @Failure 400 {object} Response "Bad Request with details"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user [post]
func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request to validator", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		userPassword, err := encryption.EncryptPassword(req.Password)
		if err != nil {
			log.Error("failed to encrypt user password", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to encrypt password"))

			return
		}

		id, err := userSaver.SaveUser(req.Username, req.Email, userPassword)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("user", req.Username))

			render.JSON(w, r, resp.Error("user already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add user"))

			return
		}

		log.Info("user saved into db", slog.Int64("id", id))

		jwtUserAccessToken, jwtUserRefreshToken, err := jwtAuth.GenerateTokens(id)
		if err != nil {
			log.Error("failed to generate json web token for user id", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to generate jwt token"))

			return
		}

		responseOK(w, r, req.Username, jwtUserAccessToken, jwtUserRefreshToken)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, username string, jwtUserAccessToken string, jwtUserRefreshToken string) {
	render.JSON(w, r, Response{
		Response:        resp.OK(),
		Username:        username,
		JWTAccessToken:  jwtUserAccessToken,
		JWTRefreshToken: jwtUserRefreshToken,
	})
}