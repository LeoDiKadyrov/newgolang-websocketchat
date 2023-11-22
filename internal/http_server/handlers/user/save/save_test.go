package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	save "new-websocket-chat/internal/http_server/handlers/user/save"
	"new-websocket-chat/internal/http_server/handlers/user/save/mocks"
	"new-websocket-chat/internal/lib/logger/handlers/slogdiscard"
	"testing"
)

func TestSaveHandler(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		email     string
		password  string
		respError string
		mockError error
	}{
		{
			name:     "Success",
			username: "AbdraBlya",
			email:    "dininchesterrr25@gmail.com",
			password: "Abdrahman_02!",
		},
		{
			name:      "Empty username",
			username:  "",
			email:     "din02winchester25@gmail.com",
			password:  "Abdrahman_02!",
			respError: "failed to decode request",
		},
		{
			name:      "Empty email",
			username:  "Abdrahman",
			email:     "",
			password:  "Abdrahman_02!",
			respError: "field email is a required field",
		},
		{
			name:      "Empty password",
			username:  "Abdrahman",
			email:     "din02winchester25@gmail.com",
			password:  "",
			respError: "field password is a required field",
		},
		{
			name:      "Empty username & email",
			username:  "",
			email:     "",
			password:  "Abdrahman_02!",
			respError: "field username is a required field",
		},
		{
			name:      "Empty email & password",
			username:  "Abdrahman",
			email:     "",
			password:  "",
			respError: "field email is a required field",
		},
		{
			name:      "Empty username & password",
			username:  "",
			email:     "din02winchester25@gmail.com",
			password:  "",
			respError: "field username is a required field",
		},
		{
			name:      "All empty",
			username:  "",
			email:     "",
			password:  "",
			respError: "field username is a required field",
		},
		{
			name:      "Invalid username",
			username:  "Aba",
			email:     "din02winchester25@gmail.com",
			password:  "Abdrahman_02!",
			respError: "field username is not a valid username",
		},
		{
			name:      "Invalid email",
			username:  "Abdrahman",
			email:     "din02winchester25",
			password:  "Abdrahman_02!",
			respError: "field email is not a valid email",
		},
		{
			name:      "Invalid password",
			username:  "Abdrahman",
			email:     "din02winchester25@gmail.com",
			password:  "abdrahman02",
			respError: "field password is not a valid password",
		},
		{
			name:      "SaveUser Error",
			username:  "AbdraBlya",
			email:     "din02winchester25@gmail.com",
			password:  "Abdrahman_02!",
			respError: "failed to add user",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userSaverMock := mocks.NewUserSaver(t)

			if test.respError == "" || test.mockError != nil {
				userSaverMock.On("SaveUser", test.username, test.email, mock.Anything). // mock.Anything because I encrypt user's password through bcrypt and hash everytime is different
													Return(int64(1), test.mockError).
													Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), userSaverMock)

			input := fmt.Sprintf(`{"username": "%s", "email": "%s", "password": "%s"}`, test.username, test.email, test.password)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var respa save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &respa))

			require.Equal(t, test.respError, respa.Error)
		})
	}
}
