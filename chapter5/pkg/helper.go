package pkg

import "golang.org/x/crypto/bcrypt"

// HashPassword genera un hash bcrypt de la contraseña
func HashPassword(password string) (string, error) {
	// GenerateFromPassword devuelve el hash bcrypt de la contraseña al costo indicado
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash ...
func CheckPasswordHash(password, hash string) bool {
	// CompareHashAndPassword compara una contraseña con hash de bcrypt
	// con su posible equivalente en texto sin formato.
	// Devuelve cero en caso de éxito o un error en caso de error.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
