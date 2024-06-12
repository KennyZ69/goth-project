package database

import (
	"fmt"
	"net/mail"
	"regexp"
)

func EmailValidator(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		fmt.Printf("The email address doesnt exist: %v", email)
		return false
	}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(email)
}

func UsernameValidator(username string) bool {
	pattern := `^[a-zA-Z0-9_-]{3,20}$`
	re := regexp.MustCompile(pattern)

	return re.MatchString(username)
}

func PasswordValidator(password string) bool {
	pattern := `^[a-zA-Z\d@$!%*?&]{8,30}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(password)
}
