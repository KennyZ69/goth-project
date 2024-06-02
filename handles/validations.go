package handles

import (
	"regexp"
	"unicode"
)

func EmailValidator(email string) bool {
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
	// pattern := `^(?:(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&]).{8,30})$`
	// re := regexp.MustCompile(pattern)
	// return re.MatchString(password)

	var (
		hasMinLength = false
		hasDigit     = false
		hasUpper     = false
		hasSpecial   = false
		minLength    = 8
	)

	if len(password) >= minLength {
		hasMinLength = true
	}

	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLength && hasDigit && hasUpper && hasSpecial
}
