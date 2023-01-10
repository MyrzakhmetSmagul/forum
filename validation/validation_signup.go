package validation

import (
	"errors"
	"regexp"
	"unicode"
)

var (
	ErrMessageValid = errors.New("Введите данные корректно")
)

func ValidationFormSignUp(uname, umail, psw, psw2 string) error {
	if uname == "" || umail == "" || psw == "" || psw2 == "" {
		return ErrMessageValid
	}
	if psw != psw2 {
		return ErrMessageValid
	}
	if !checkPatternName(uname) {
		return ErrMessageValid
	}
	if !checkPatternEmail(umail) {
		return ErrMessageValid
	}
	if !isValidPassword(psw) {
		return ErrMessageValid
	}
	return nil
}

func checkPatternName(uname string) bool {
	var re = regexp.MustCompile("^[a-zA-Z0-9]+$")
	if re.MatchString(uname) {
		return true
	}
	return false
}

func checkPatternEmail(umail string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(umail) {
		return true
	}
	return false
}

func isValidPassword(psw string) bool {
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
