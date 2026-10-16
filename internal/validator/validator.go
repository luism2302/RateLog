package validator

import "regexp"

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) AddError(key string, value string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = value
	}
}

func (v *Validator) GetError(key string) string {
	return v.Errors[key]
}

func (v *Validator) Check(ok bool, key, value string) {
	if !ok {
		v.AddError(key, value)
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func MatchesRX(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func ValidateEmail(email string, v *Validator) {
	v.Check(email != "", "email", "Cannot be blank")
	v.Check(MatchesRX(email, EmailRX), "email", "Must be a valid email address")
}

func ValidatePassword(password string, v *Validator) {
	v.Check(password != "", "password", "Cannot be blank")
	v.Check(len(password) >= 8, "password", "Must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "Must not be more than 72 bytes long")
}
