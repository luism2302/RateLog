package views

import (
	"time"

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
	validator.ValidatePassword(values.Password, v)
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
	v.Check(len(values.ConfirmPassword) >= 8, "confirmPassword", "Must be at least 8 bytes long")
	v.Check(len(values.ConfirmPassword) <= 72, "confirmPassword", "Must not be more than 72 bytes long")
	v.Check(values.Password == values.ConfirmPassword, "confirmPassword", "Passwords don't match")
	v.Check(values.Name != "", "name", "Cannot be blank")
	v.Check(len(values.Name) >= 5, "name", "Must be at least 5 bytes long")
	v.Check(len(values.Name) <= 15, "name", "Must not be more than 15 bytes long")
}

type AddGameFormValues struct {
	ID              int
	Name            string
	Released        string
	BackgroundImage string
	Added           int
}

func NewAddGameFormValues(id, added int, name, released, backgroundImage string) AddGameFormValues {
	return AddGameFormValues{
		ID:              id,
		Name:            name,
		Released:        released,
		BackgroundImage: backgroundImage,
		Added:           added,
	}
}

func ValidateAddGameForm(values AddGameFormValues, v *validator.Validator) {
	v.Check(values.Name != "", "name", "cannot be blank")
	v.Check(values.Released != "", "released", "cannot be blank")
	v.Check(values.BackgroundImage != "", "backgroundImage", "cannot be blank")

	_, err := time.Parse("2006-01-02 15:04:05 -0700 MST", values.Released)
	if err != nil {
		v.AddError("released", "must be a valid date")
	}
}
