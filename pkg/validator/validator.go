package validator

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("positive", validatePositive)
		v.RegisterValidation("non_negative", validateNonNegative)
	}
}

func validatePositive(fl validator.FieldLevel) bool {
	switch v := fl.Field().Interface().(type) {
	case int, int8, int16, int32, int64:
		return fl.Field().Int() > 0
	case uint, uint8, uint16, uint32, uint64:
		return fl.Field().Uint() > 0
	case float32, float64:
		return fl.Field().Float() > 0
	default:
		_ = v
		return false
	}
}

func validateNonNegative(fl validator.FieldLevel) bool {
	switch v := fl.Field().Interface().(type) {
	case int, int8, int16, int32, int64:
		return fl.Field().Int() >= 0
	case uint, uint8, uint16, uint32, uint64:
		return true // unsigned integers are always non-negative
	case float32, float64:
		return fl.Field().Float() >= 0
	default:
		_ = v
		return false
	}
}

func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			messages = append(messages, fmt.Sprintf("Field '%s' %s", e.Field(), getErrorMessage(e)))
		}
		return strings.Join(messages, "; ")
	}
	return err.Error()
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", e.Param())
	case "positive":
		return "must be positive"
	case "non_negative":
		return "must be non-negative"
	default:
		return fmt.Sprintf("failed validation '%s'", e.Tag())
	}
}
