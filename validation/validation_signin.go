package validation

import (
	"regexp"
	"unicode"
)

func ValidationFormSignIn(umail, psw string) error {
	if umail == "" || psw == "" {
		return ErrMessageValid
	}
	if !checkPatternEmailForSignIn(umail) {
		return ErrMessageValid

	}
	if !isValidPasswordForSignIn(psw) {
		return ErrMessageValid

	}
	return nil
}

func checkPatternEmailForSignIn(umail string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(umail) {
		return true
	}
	return false
}

func isValidPasswordForSignIn(psw string) bool {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)
	if len(psw) >= 7 {
		hasMinLen = true
	}
	for _, char := range psw {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber
}
