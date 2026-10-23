package views

import (
	"github.com/luism2302/RateLog/internal/validator"
)

type LoginFormValues struct {
	Email    string
	Password string
}

func NewLoginFormValues(email, password string) LoginFormValues {
	return LoginFormValues{Email: email, Password: password}
}

func ValidateLoginForm(values LoginFormValues, v *validator.Validator) {
	validator.ValidateEmail(values.Email, v)
	v.Check(values.Password != "", "password", "Cannot be blank")
}

type SignupFormValues struct {
	Email           string
	Name            string
	Password        string
	ConfirmPassword string
}

func NewSignupFormValues(email, name, password, confirmPassword string) SignupFormValues {
	return SignupFormValues{Email: email, Name: name, Password: password, ConfirmPassword: confirmPassword}
}

func ValidateSignupForm(values SignupFormValues, v *validator.Validator) {
	validator.ValidateEmail(values.Email, v)
	validator.ValidatePassword(values.Password, v)
	v.Check(values.ConfirmPassword != "", "confirmPassword", "Cannot be blank")
	v.Check(len(values.ConfirmPassword) >= 8, "confirmPassword", "Must be at least 8 characters long")
	v.Check(len(values.ConfirmPassword) <= 72, "confirmPassword", "Must not be more than 72 characters long")
	v.Check(values.Password == values.ConfirmPassword, "confirmPassword", "Passwords don't match")
	v.Check(values.Name != "", "name", "Cannot be blank")
	v.Check(len(values.Name) >= 5, "name", "Must be at least 5 characters long")
	v.Check(len(values.Name) <= 15, "name", "Must not be more than 15 characters long")
}

type ReviewFormValues struct {
	Body   string
	Rating int
}

func NewReviewFormValues(body string, rating int) ReviewFormValues {
	return ReviewFormValues{
		Body:   body,
		Rating: rating,
	}
}

func ValidateReviewGameForm(values ReviewFormValues, v *validator.Validator) {
	v.Check(values.Body != "", "body", "Cannot be blank")
	v.Check(values.Rating >= 1, "rating", "Cannot be lower than 0")
	v.Check(values.Rating <= 5, "rating", "Cannot be greater than 5")
}
