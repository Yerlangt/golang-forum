package service

import "strings"

func validateEmail(email string) (bool, string) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, "Email without domian"
	}
	check2 := strings.Split(parts[1], ".")
	if len(check2) < 2 {
		return false, "Imporoper email domain"
	}
	for i := range check2 {
		if len(check2[i]) < 2 || len(check2[i]) > 253 {
			return false, "Imporoper email domain"
		}
	}
	if len(parts[0]) < 2 || len(parts[0]) > 64 {
		return false, "Imporoper email recipient"
	}
	for i := range parts[0] {
		if parts[0][i] != '!' && parts[0][i] != '#' && parts[0][i] != '$' && parts[0][i] != '%' && parts[0][i] != '&' && parts[0][i] != '\'' &&
			parts[0][i] != '*' && parts[0][i] != '+' && parts[0][i] != '-' && parts[0][i] != '/' && parts[0][i] != '=' && parts[0][i] != '?' &&
			parts[0][i] != '^' && parts[0][i] != '_' && parts[0][i] != '`' && parts[0][i] != '|' && parts[0][i] != '.' &&
			(parts[0][i] < 'a' || parts[0][i] > 'z') && (parts[0][i] < 'A' || parts[0][i] > 'Z') && (parts[0][i] < '0' || parts[0][i] > '9') {
			return false, "Imporoper email recipient"
		}
	}
	return true, ""
}

func validatePassword(pswd string) (bool, string) {
	if len(pswd) < 5 {
		return false, "Password too short"
	}
	for i := range pswd {
		if pswd[i] != '@' && pswd[i] != '#' && pswd[i] != '&' && pswd[i] != '!' && pswd[i] != '$' && pswd[i] != '%' && pswd[i] != '*' &&
			(pswd[i] < 'a' || pswd[i] > 'z') && (pswd[i] < 'A' || pswd[i] > 'Z') && (pswd[i] < '0' || pswd[i] > '9') {
			return false, "Improper characters in password"
		}
	}
	return true, ""
}

func validateUsername(name string) (bool, string) {
	if len(name) < 5 && len(name) > 15 {
		return false, "Username too short"
	}
	for i := range name {
		if (name[i] < 'a' || name[i] > 'z') && (name[i] < 'A' || name[i] > 'Z') && (name[i] < '0' || name[i] > '9') {
			return false, "Improper characters in username"
		}
	}
	return true, ""
}
