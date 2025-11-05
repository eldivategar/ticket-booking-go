package validator

import (
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	v := validator.New()

	// Register custom validation here
	v.RegisterValidation("base64", validateBase64)

	return &Validator{
		validate: v,
	}
}

func (v *Validator) Validate(data interface{}) error {
	return v.validate.Struct(data)
}

// FormatErrors processes validation errors into a neat map.
func (v *Validator) FormatErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrs {
			field := strings.ToLower(fe.Field())
			errors[field] = v.formatErrorMessage(fe)
		}
	}

	return errors
}

// formatErrorMessage generates user-friendly error messages based on validation tags.
func (v *Validator) formatErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", fe.Param())
	case "base64":
		return "must be a valid base64 encoded string"
	default:
		// Fallback message
		return fmt.Sprintf("invalid value for tag '%s'", fe.Tag())
	}
}

// validateBase64 is a custom validator function to check if a string is valid base64.
func validateBase64(fl validator.FieldLevel) bool {
	// Ambil value dari field
	str := fl.Field().String()

	// Jika string kosong, anggap valid (biarkan tag 'required' yang menangani)
	if str == "" {
		return true
	}

	// Cek apakah ini data URI (memiliki prefix)
	if strings.Contains(str, ",") {
		parts := strings.SplitN(str, ",", 2)
		if len(parts) != 2 {
			return false // Format data URI salah
		}
		// Ambil hanya bagian data-nya saja
		str = parts[1]
	}

	// Coba decode string base64
	_, err := base64.StdEncoding.DecodeString(str)

	// Jika tidak ada error (err == nil), berarti valid
	return err == nil
}

// Password must be at least 8 chars, with uppercase, lowercase, digit, and special char
func password(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
