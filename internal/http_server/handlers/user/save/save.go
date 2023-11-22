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
	"new-websocket-chat/internal/lib/logger/sl"
	"new-websocket-chat/internal/storage"
)

type Request struct {
	Username string `json:"username" validate:"required,min=4,max=24"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20,containsany=!@#?"` // validation, will it return error if password empty or rules don't apply???
}

type Response struct {
	resp.Response
	Username string `json:"username,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.37.1 --name=UserSaver
type UserSaver interface {
	SaveUser(username string, email string, password string) (int64, error)
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		log.Info("r.Body is following", slog.Any("req", req))

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

		id, err := userSaver.SaveUser(req.Username, req.Email, req.Password) // TODO: add password hashing before saving
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
		responseOK(w, r, req.Username)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, username string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Username: username,
	})
}
