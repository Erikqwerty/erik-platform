package validator

import (
	"errors"
	"fmt"
	"regexp"
)

// ValidEmails - возвращает ошибку в виде списока не валидных email
func ValidEmails(emails []string) error {
	var errtext string

	for _, email := range emails {
		if !IsValidEmail(email) {
			errtext += fmt.Sprintf("email: %v не валиден;\n", email)
		}
	}

	if errtext != "" {
		return errors.New(errtext)
	}

	return nil
}

// IsValidEmail проверяет валидность email-адреса. Возвращает true если валидно.
func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}
