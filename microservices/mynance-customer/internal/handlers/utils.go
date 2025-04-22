package handlers

import "regexp"

func isValidPassword(password string) bool {
	var (
		upperCase = regexp.MustCompile(`[A-Z]`)
		lowerCase = regexp.MustCompile(`[a-z]`)
		digit     = regexp.MustCompile(`\d`)
		special   = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	)
	return len(password) >= 8 && upperCase.MatchString(password) && lowerCase.MatchString(password) && digit.MatchString(password) && special.MatchString(password)
}
