package encryption

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}
