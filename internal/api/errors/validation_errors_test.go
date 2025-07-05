package errors_test

import (
	"testing"

	"web-scraper.dev/internal/api/errors"
	"web-scraper.dev/internal/utils/validator"
)

func TestToValidationErrors(t *testing.T) {
	t.Parallel()

	v := validator.New()
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "required",
			input: struct {
				Email string `json:"email" validate:"required"`
			}{},
			expected: "email is a required field",
		},
		{
			name: "eqfield",
			input: struct {
				Password        string `json:"password"`
				ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
			}{
				Password:        "abc123",
				ConfirmPassword: "abc213",
			},
			expected: "confirm_password must be equal to Password",
		},
		{
			name: "min",
			input: struct {
				Password string `json:"password" validate:"min=8"`
			}{Password: "123456"},
			expected: "password must be at least in 8 characters",
		},
		{
			name: "email",
			input: struct {
				Email string `json:"email" validate:"email"`
			}{Email: "mail@"},
			expected: "email must be a valid email address",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := v.Struct(tc.input)
			if errResp := errors.ToValidationErrors(err); errResp == nil || len(errResp.Errors) != 1 {
				t.Fatalf(`Expected:"{[%v]}", Got:"%v"`, tc.expected, errResp)
			} else if errResp.Errors[0] != tc.expected {
				t.Fatalf(`Expected:"%v", Got:"%v"`, tc.expected, errResp.Errors[0])
			}
		})
	}
}
