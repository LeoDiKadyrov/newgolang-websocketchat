package refresh_test

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	refresh "new-websocket-chat/internal/http_server/handlers/jwt"
	"new-websocket-chat/internal/http_server/handlers/jwt/mocks"
	"new-websocket-chat/internal/lib/logger/handlers/slogdiscard"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestRefreshHandler(t *testing.T) {
	tests := []struct {
		name           string
		tokenFunc      func(int64) (string, error)
		expectError    bool
		expectedStatus int
		setupMock      func(*mocks.TokenService, string, int64)
	}{
		{
			name:           "Success",
			tokenFunc:      generateTestRefreshToken,
			expectError:    false,
			expectedStatus: http.StatusOK,
			setupMock: func(m *mocks.TokenService, token string, userID int64) {
				m.On("ExtractToken", mock.Anything).Return(token, nil)
				m.On("ValidateToken", token).Return(&jwt.StandardClaims{Subject: strconv.FormatInt(userID, 10)}, nil)
				m.On("GenerateTokens", userID).Return("new_access_token", "new_refresh_token", nil)
			},
		},
		{
			name:           "ExpiredToken",
			tokenFunc:      generateExpiredTestRefreshToken,
			expectError:    true,
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(m *mocks.TokenService, token string, userID int64) {
				m.On("ExtractToken", mock.Anything).Return(token, nil)
				m.On("ValidateToken", token).Return(nil, fmt.Errorf("invalid or expired token"))
			},
		},
		{
			name:           "Token Extraction Failure",
			tokenFunc:      generateTestRefreshToken,
			expectError:    true,
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(m *mocks.TokenService, token string, userID int64) {
				m.On("ExtractToken", mock.Anything).Return("", fmt.Errorf("failed to extract token"))
			},
		},
		{
			name:           "Token Validation Error",
			tokenFunc:      generateTestRefreshToken,
			expectError:    true,
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(m *mocks.TokenService, token string, userID int64) {
				m.On("ExtractToken", mock.Anything).Return(token, nil)
				m.On("ValidateToken", token).Return(nil, fmt.Errorf("invalid token"))
			},
		},
		{
			name:           "User ID Parsing Error",
			tokenFunc:      generateTestRefreshToken,
			expectError:    true,
			expectedStatus: http.StatusBadRequest,
			setupMock: func(m *mocks.TokenService, token string, userID int64) {
				m.On("ExtractToken", mock.Anything).Return(token, nil)
				m.On("ValidateToken", token).Return(&jwt.StandardClaims{Subject: "invalid"}, nil)
			},
		},
		{
			name:           "Token Generation Error",
			tokenFunc:      generateTestRefreshToken,
			expectError:    true,
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(m *mocks.TokenService, token string, userID int64) {
				m.On("ExtractToken", mock.Anything).Return(token, nil)
				m.On("ValidateToken", token).Return(&jwt.StandardClaims{Subject: strconv.FormatInt(userID, 10)}, nil)
				m.On("GenerateTokens", userID).Return("", "", fmt.Errorf("token generation failed"))
			},
		},
	}

	testUserID := int64(12345)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testToken, err := test.tokenFunc(testUserID)
			require.NoError(t, err)

			tokenServiceMock := mocks.NewTokenService(t)
			test.setupMock(tokenServiceMock, testToken, testUserID)

			handler := refresh.New(slogdiscard.NewDiscardLogger(), tokenServiceMock)

			req, err := http.NewRequest(http.MethodPost, "/api/jwt/refresh", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+testToken)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if test.expectError {
				require.NotEqual(t, test.expectedStatus, rr.Code)
			} else {
				require.Equal(t, rr.Code, http.StatusOK)
				var resp refresh.Response
				err = json.Unmarshal(rr.Body.Bytes(), &resp)
				require.NoError(t, err, "Error unmarshalling response body")
				require.NotEmpty(t, resp.JWTAccessToken, "JWTAccessToken should not be empty")
				require.NotEmpty(t, resp.JWTRefreshToken, "JWTRefreshToken should not be empty")
			}
		})
	}
}

func generateTestRefreshToken(userID int64) (testRefreshToken string, err error) {
	const op = "internal.http_server.handlers.jwt.generateTestToken"

	jwtIssuer := os.Getenv("JWT_ISSUER")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	refreshTokenClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    jwtIssuer,
		Subject:   strconv.FormatInt(userID, 10),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	testRefreshToken, err = refreshToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign refresh token: %w", op, err)
	}

	return testRefreshToken, nil
}

func generateExpiredTestRefreshToken(userID int64) (testRefreshToken string, err error) {
	const op = "internal.http_server.handlers.jwt.generateExpiredTestRefreshToken"

	jwtIssuer := os.Getenv("JWT_ISSUER")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	expiredTime := time.Now().Add(-24 * time.Hour)
	refreshTokenClaims := jwt.StandardClaims{
		ExpiresAt: expiredTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    jwtIssuer,
		Subject:   strconv.FormatInt(userID, 10),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	testRefreshToken, err = refreshToken.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("%s: failed to sign refresh token: %w", op, err)
	}

	return testRefreshToken, nil
}
