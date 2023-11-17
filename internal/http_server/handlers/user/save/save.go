package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "new-websocket-chat/internal/lib/api/response"
	"new-websocket-chat/internal/lib/logger/sl"
)

type UserSaver interface {
	SaveUser(username string, email string, password_hash string) (int64, error)
}

type Request struct {
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,password"` // validation
}

type Response struct {
	resp.Response
	Username string `json:"username,omitempty"`
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decore request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request to validator", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

	}
}
