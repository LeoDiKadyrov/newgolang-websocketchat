package encryption

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestNewEncryptedPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "password = Abdrahman",
			password: "Abdrahman",
		},
		{
			name:     "empty password",
			password: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hashedPassword, err := EncryptPassword(test.password)
			if err != nil {
				t.Errorf("EncryptPassword(%q) returned an error %s", test.password, err)
			}
			if hashedPassword == "" {
				t.Errorf("EncryptPassword(%q) returned an empty string", test.password)
			}

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(test.password))
			if err != nil {
				t.Errorf("Failed to verify password '%s' with its hash '%s': %s", test.password, hashedPassword, err)
			}
		})
	}
}
