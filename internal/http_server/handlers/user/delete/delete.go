package delete

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

type DeleteRequest struct {
	Username string `json:"username" validate:"required,min=4,max=24"`
	Email    string `json:"email" validate:"required,email"`
}

type Response struct {
	resp.Response
	Username string `json:"username,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.37.1 --name=UserDeleter
type UserDeleter interface {
	DeleteUser(username string, email string) error
}

func New(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.delete.new"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req DeleteRequest

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

		err = userDeleter.DeleteUser(req.Username, req.Email)
		if errors.Is(err, storage.ErrUsernameNotFound) || errors.Is(err, storage.ErrEmailNotFound) {
			log.Info("username or email not found", slog.String("user", req.Username))

			render.JSON(w, r, resp.Error("user not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete user"))

			return
		}

		log.Info("user deleted from db", slog.Any("err", err))
		responseOK(w, r, req.Username)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, username string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Username: username,
	})
}

/* Clean code thoughts & questions to myself
TODO:
[ ] Validation depends on go-playground/validator/v10 - same comment as in user.save
*/
