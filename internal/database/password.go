package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// GenerateFromPassword creates a hashed password with a salt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	// CompareHashAndPassword compares the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func mainFunc() {
	password := "your_password_here"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Now you would store `hashedPassword` in your database
	log.Println("Hashed Password:", hashedPassword)

	isMatch := CheckPasswordHash(password, hashedPassword)
	if isMatch {
		log.Println("Password is correct")
	} else {
		log.Println("Password is incorrect")
	}
}
