package delete_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	delete "new-websocket-chat/internal/http_server/handlers/user/delete"
	"new-websocket-chat/internal/http_server/handlers/user/delete/mocks"
	"new-websocket-chat/internal/lib/logger/handlers/slogdiscard"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		email     string
		respError string
		mockError error
	}{
		{
			name:     "Success",
			username: "AbdraBlya",
			email:    "dininchesterrr25@gmail.com",
		},
		{
			name:      "Not exist",
			username:  "AbdraBlyaaaaa",
			email:     "dininchesterrdfsdfr25@gmail.com",
			respError: "user not found",
		},
		{
			name:      "Empty username",
			username:  "",
			email:     "din02winchester25@gmail.com",
			respError: "failed to decode request",
		},
		{
			name:      "Empty email",
			username:  "Abdrahman",
			email:     "",
			respError: "field email is a required field",
		},
		{
			name:      "Empty username & email",
			username:  "",
			email:     "",
			respError: "field username is a required field",
		},
		{
			name:      "Invalid username",
			username:  "Aba",
			email:     "din02winchester25@gmail.com",
			respError: "field username is not a valid username",
		},
		{
			name:      "Invalid email",
			username:  "Abdrahman",
			email:     "din02winchester25",
			respError: "field email is not a valid email",
		},
		{
			name:      "DeleteUser Error",
			username:  "Abdrahmanishe",
			email:     "din02winchester25@gmail.com",
			respError: "failed to delete user",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userDeleterMock := mocks.NewUserDeleter(t)

			if test.respError == "" || test.mockError != nil {
				userDeleterMock.On("DeleteUser", test.username, test.email).
					Return(test.mockError).
					Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), userDeleterMock)

			input := fmt.Sprintf(`{"username": "%s", "email": "%s"}`, test.username, test.email)

			req, err := http.NewRequest(http.MethodPost, "/delete", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp delete.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, test.respError, resp.Error)
		})
	}
}
