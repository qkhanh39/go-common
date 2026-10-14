package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/qkhanh39/go-common/apperrors"
)

type FieldError = validator.FieldError

type Validator struct {
	validate       *validator.Validate
	customCodes    map[string]string
	customMessages map[string]func(fe FieldError) string
}

func New() *Validator {
	v := &Validator{
		validate:       validator.New(validator.WithRequiredStructEnabled()),
		customCodes:    make(map[string]string),
		customMessages: make(map[string]func(fe FieldError) string),
	}

	v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

func (v *Validator) Validate(s any) error {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		ve := apperrors.NewValidationError()
		for _, fe := range valErrs {
			ve.Details = append(ve.Details, apperrors.ErrorDetail{
				Field: fe.Field(),
				Code:  v.getTagCode(fe.Tag(), fe.Param()),
				Issue: v.getIssueMessage(fe),
			})
		}
		return ve
	}

	return apperrors.InternalError{Issue: "validation failed"}
}

func (v *Validator) Underlying() *validator.Validate {
	return v.validate
}

func (v *Validator) RegisterValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

func (v *Validator) RegisterTagCode(tag, code string) *Validator {
	v.customCodes[tag] = code
	return v
}

func (v *Validator) RegisterTagMessage(tag string, fn func(fe FieldError) string) *Validator {
	v.customMessages[tag] = fn
	return v
}

func (v *Validator) getIssueMessage(fe FieldError) string {
	if fn, ok := v.customMessages[fe.Tag()]; ok {
		return fn(fe)
	}
	return defaultIssueMessage(fe)
}

func (v *Validator) getTagCode(tag, param string) string {
	if code, ok := v.customCodes[tag]; ok {
		return code
	}
	return defaultTagCode(tag, param)
}

func defaultTagCode(tag, param string) string {
	switch tag {
	case "required":
		return "REQUIRED"
	case "email":
		return "INVALID_EMAIL"
	case "min":
		return "MIN_LENGTH"
	case "max":
		return "MAX_LENGTH"
	case "len":
		return "EXACT_LENGTH"
	case "uuid", "uuid4":
		return "INVALID_UUID"
	case "url":
		return "INVALID_URL"
	case "eqfield":
		return "MISMATCH"
	case "oneof":
		return "INVALID_OPTION"
	case "gt":
		return "GREATER_THAN"
	case "gte":
		return "GREATER_THAN_OR_EQUAL"
	case "lt":
		return "LESS_THAN"
	case "lte":
		return "LESS_THAN_OR_EQUAL"
	case "numeric":
		return "MUST_BE_NUMERIC"
	case "alpha":
		return "MUST_BE_ALPHA"
	case "alphanum":
		return "MUST_BE_ALPHANUMERIC"
	default:
		return strings.ToUpper(tag)
	}
}

// defaultIssueMessage creates a human-readable error message.
func defaultIssueMessage(fe FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "len":
		return "must be exactly " + fe.Param() + " characters"
	case "uuid", "uuid4":
		return "must be a valid UUID"
	case "url":
		return "must be a valid URL"
	case "eqfield":
		return "must match " + fe.Param()
	case "oneof":
		return "must be one of: " + fe.Param()
	case "gt":
		return "must be greater than " + fe.Param()
	case "gte":
		return "must be greater than or equal to " + fe.Param()
	case "lt":
		return "must be less than " + fe.Param()
	case "lte":
		return "must be less than or equal to " + fe.Param()
	case "numeric":
		return "must contain only numbers"
	case "alpha":
		return "must contain only letters"
	case "alphanum":
		return "must contain only letters and numbers"
	default:
		return fe.Error()
	}
}
