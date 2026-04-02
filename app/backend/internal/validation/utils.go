package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/apk471/go-boilerplate/internal/errs"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validatable interface {
	Validate() error
}

type CustomValidationError struct {
	Field   string
	Message string
}

type CustomValidationErrors []CustomValidationError

func (c CustomValidationErrors) Error() string {
	return "Validation failed"
}

func BindAndValidate(c echo.Context, payload Validatable) error {
	if err := bindRequest(c, payload); err != nil {
		return err
	}

	if msg, fieldErrors := validateStruct(payload); fieldErrors != nil {
		return errs.NewBadRequestError(msg, true, nil, fieldErrors, nil)
	}

	return nil
}

func validateStruct(v Validatable) (string, []errs.FieldError) {
	if err := v.Validate(); err != nil {
		return extractValidationErrors(err)
	}
	return "", nil
}

func extractValidationErrors(err error) (string, []errs.FieldError) {
	var fieldErrors []errs.FieldError
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		customValidationErrors, ok := err.(CustomValidationErrors)
		if !ok {
			return "Validation failed", []errs.FieldError{{
				Field: "body",
				Error: err.Error(),
			}}
		}
		for _, err := range customValidationErrors {
			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: err.Field,
				Error: err.Message,
			})
		}

		return "Validation failed", fieldErrors
	}

	for _, err := range validationErrors {
		field := strings.ToLower(err.Field())
		var msg string

		switch err.Tag() {
		case "required":
			msg = "is required"
		case "min":
			if err.Type().Kind() == reflect.String {
				msg = fmt.Sprintf("must be at least %s characters", err.Param())
			} else {
				msg = fmt.Sprintf("must be at least %s", err.Param())
			}
		case "max":
			if err.Type().Kind() == reflect.String {
				msg = fmt.Sprintf("must not exceed %s characters", err.Param())
			} else {
				msg = fmt.Sprintf("must not exceed %s", err.Param())
			}
		case "oneof":
			msg = fmt.Sprintf("must be one of: %s", err.Param())
		case "email":
			msg = "must be a valid email address"
		case "e164":
			msg = "must be a valid phone number with country code"
		case "uuid":
			msg = "must be a valid UUID"
		case "uuidList":
			msg = "must be a comma-separated list of valid UUIDs"
		case "dive":
			msg = "some items are invalid"
		default:
			if err.Param() != "" {
				msg = fmt.Sprintf("%s: %s:%s", field, err.Tag(), err.Param())
			} else {
				msg = fmt.Sprintf("%s: %s", field, err.Tag())
			}
		}

		fieldErrors = append(fieldErrors, errs.FieldError{
			Field: strings.ToLower(err.Field()),
			Error: msg,
		})
	}

	return "Validation failed", fieldErrors
}

func bindRequest(c echo.Context, payload Validatable) error {
	binder := &echo.DefaultBinder{}

	if err := binder.BindPathParams(c, payload); err != nil {
		return errs.NewBadRequestError("Invalid path parameters", true, nil, []errs.FieldError{
			{Field: "path", Error: "contains invalid values"},
		}, nil)
	}

	if method := c.Request().Method; method == http.MethodGet || method == http.MethodDelete {
		if err := binder.BindQueryParams(c, payload); err != nil {
			return errs.NewBadRequestError("Invalid query parameters", true, nil, []errs.FieldError{
				{Field: "query", Error: "contains invalid values"},
			}, nil)
		}
	}

	if shouldBindBody(c.Request().Method) {
		if err := bindJSONBody(c, payload); err != nil {
			return err
		}
	}

	return nil
}

func shouldBindBody(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}

func bindJSONBody(c echo.Context, payload Validatable) error {
	req := c.Request()
	if req.Body == nil {
		return nil
	}

	if req.ContentLength == 0 {
		return nil
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(payload); err != nil {
		return mapJSONDecodeError(err)
	}

	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		return errs.NewBadRequestError("Invalid JSON body", true, nil, []errs.FieldError{
			{Field: "body", Error: "must contain a single JSON object"},
		}, nil)
	}

	return nil
}

func mapJSONDecodeError(err error) error {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return errs.NewBadRequestError("Invalid JSON body", true, nil, []errs.FieldError{
			{Field: "body", Error: "contains malformed JSON"},
		}, nil)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := strings.TrimSpace(typeErr.Field)
		if field == "" {
			field = "body"
		}
		return errs.NewBadRequestError("Invalid JSON body", true, nil, []errs.FieldError{
			{Field: strings.ToLower(field), Error: fmt.Sprintf("must be a valid %s", typeErr.Type.String())},
		}, nil)
	}

	if err == io.EOF {
		return nil
	}

	const unknownFieldPrefix = "json: unknown field "
	if strings.HasPrefix(err.Error(), unknownFieldPrefix) {
		field := strings.Trim(err.Error()[len(unknownFieldPrefix):], "\"")
		return errs.NewBadRequestError("Invalid JSON body", true, nil, []errs.FieldError{
			{Field: field, Error: "is not allowed"},
		}, nil)
	}

	return errs.NewBadRequestError("Invalid JSON body", true, nil, []errs.FieldError{
		{Field: "body", Error: "must be a valid JSON object"},
	}, nil)
}

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func IsValidUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}
