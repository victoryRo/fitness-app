package internal

import "golang.org/x/crypto/bcrypt"

// HashPassword receives a password string and returns a hash password
func HashPassword(password string) string {
	bytes, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return ""
	}
	return string(bytes)
}

// CheckPasswordHash compares the user's password string with the database password hash
func CheckPasswordHash(password, hash string) bool {
	e := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return e == nil
}
