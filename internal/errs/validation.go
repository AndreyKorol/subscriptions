package errs

import (
	"errors"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type ValidationDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FromValidator(err error) []ValidationDetail {
	vErrors, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		return nil
	}

	details := make([]ValidationDetail, 0, len(vErrors))
	for _, fe := range vErrors {
		details = append(details, ValidationDetail{
			Field:   fieldName(fe),
			Message: tagMessage(fe.Tag(), fe.Param()),
		})
	}
	return details
}

func fieldName(fe validator.FieldError) string {
	name := fe.StructField()
	return toSnake(name)
}

func toSnake(s string) string {
	var buf []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			buf = append(buf, '_')
		}
		buf = append(buf, unicode.ToLower(r))
	}
	return string(buf)
}

func tagMessage(tag, param string) string {
	switch tag {
	case "required":
		return "field is required"
	case "ne":
		return "field must not be empty"
	case "gte":
		if param == "0" {
			return "must be 0 or greater"
		}
		return "must be " + param + " or greater"
	case "lte":
		return "must be " + param + " or less"
	case "uuid4":
		return "must be a valid UUID"
	case "datetime":
		return "must be in format " + param
	case "min":
		return "minimum is " + param
	case "max":
		return "maximum is " + param
	default:
		return "validation failed on " + tag
	}
}


