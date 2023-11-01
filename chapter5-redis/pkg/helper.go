package pkg

import "golang.org/x/crypto/bcrypt"

// HashPassword HahsPassword converts a string to a hashed password
// returns a string and erorr
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks if a string is equal to its hash password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
