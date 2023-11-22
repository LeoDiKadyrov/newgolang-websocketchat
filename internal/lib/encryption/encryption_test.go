package encryption

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

func TestNewEncryptedPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "only english letters",
			password:    "Abdrahman",
			expectError: false,
		},
		{
			name:        "empty password",
			password:    "",
			expectError: false,
		},
		{
			name:        "with special symbols",
			password:    "Abdrahman_",
			expectError: false,
		},
		{
			name:        "with numbers",
			password:    "Abdrahman02",
			expectError: false,
		},
		{
			name:        "with numbers and symbols",
			password:    "Abdrahman_02",
			expectError: false,
		},
		{
			name:        "long password",
			password:    strings.Repeat("a", 72), // 72 is bcrypt limit
			expectError: false,
		},
		{
			name:        "extremely long password",
			password:    strings.Repeat("a", 73), // 72 is bcrypt limit - should return error
			expectError: true,
		},
		{
			name:        "only special symbols",
			password:    "!@#$%^&*()_+", // 72 is bcrypt limit - should return error
			expectError: false,
		},
		{
			name:        "only numbers",
			password:    "123456789", // 72 is bcrypt limit - should return error
			expectError: false,
		},
		{
			name:        "only numbers and symbols",
			password:    "!@#$%^&*()_+123456789", // 72 is bcrypt limit - should return error
			expectError: false,
		},
		{
			name:        "uppercase only",
			password:    "ABDRAHMAN", // 72 is bcrypt limit - should return error
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hashedPassword, err := EncryptPassword(test.password)
			if (err != nil) != test.expectError {
				t.Errorf("EncryptPassword(%q) returned an error %s", test.password, err)
			}
			if !test.expectError {
				if hashedPassword == "" {
					t.Errorf("EncryptPassword(%q) returned an empty string", test.password)
				}

				err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(test.password))
				if err != nil {
					t.Errorf("Failed to verify password '%s' with its hash '%s': %s", test.password, hashedPassword, err)
				}
			}
		})
	}
}
