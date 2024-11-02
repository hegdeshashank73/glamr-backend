package utils

import "regexp"

var USERNAME_REGEX *regexp.Regexp
var EMAIL_REGEX *regexp.Regexp

func init() {
	EMAIL_REGEX = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	USERNAME_REGEX = regexp.MustCompile(`^[a-z]+[a-z0-9_]+$`)
}

func ValidateEmail(email string) bool {
	return EMAIL_REGEX.MatchString(email)
}
