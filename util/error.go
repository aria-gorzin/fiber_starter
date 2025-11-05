package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// CustomError struct to standardize error responses
type CustomError struct {
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e CustomError) Error() string {
	return fmt.Sprintf("ErrorCode: %s, Message: %s", e.ErrorCode, e.Message)
}

// ValidateError struct to standardize validation error responses
type ValidateError struct {
	ErrorCode string      `json:"error_code"`
	Message   string      `json:"message"`
	Field     string      `json:"field"`
	Tag       string      `json:"tag"`
	Value     interface{} `json:"value"`
}

type ValidateErrors []ValidateError

func (e ValidateErrors) Error() string {
	return "validation error"
}

var validationMessages = map[string]string{
	"required": "%s is required.",
	"email":    "%s must be a valid email address.",
	"min":      "%s must be at least %s.",
	"max":      "%s must be at most %s.",
	"oneof":    "%s must be one of the following: %s.",
}

func NewValidationErrors(errs validator.ValidationErrors) ValidateErrors {
	validationErrors := []ValidateError{}
	for _, err := range errs {
		validationErrors = append(validationErrors, NewValidateError(err))
	}

	return validationErrors
}

func NewValidateError(err validator.FieldError) ValidateError {
	template := validationMessages[err.Tag()]
	if template == "" {
		template = "%s is invalid."
	}
	var message string
	switch err.Tag() {
	case "required", "email", "url":
		message = fmt.Sprintf(template, err.Field())

	case "min", "max":
		message = fmt.Sprintf(template, err.Field(), err.Param())

	case "datetime":
		message = fmt.Sprintf(template, err.Field(), err.Param())

	case "oneof":
		paramFormatted := strings.ReplaceAll(err.Param(), " ", ", ")
		message = fmt.Sprintf(template, err.Field(), paramFormatted)

	default:
		message = fmt.Sprintf("%s failed validation on '%s'", err.Field(), err.Tag())
	}

	return ValidateError{
		ErrorCode: "E3001",
		Message:   message,
		Field:     err.Field(),
		Tag:       err.Tag(),
		Value:     err.Value(),
	}
}

// ErrorResponse function to send a standardized error response
func ErrorResponse(c *fiber.Ctx, err CustomError) error {
	return c.Status(err.StatusCode).JSON(err)
}

// ValidationErrorResponse function to send a standardized validation error response
func ValidationErrorResponse(c *fiber.Ctx, errs ValidateErrors) error {
	return c.Status(fiber.StatusBadRequest).JSON(errs)
}

// NewConflictError function to create a new conflict error
func NewConflictError(table, column, value string) ValidateErrors {
	return []ValidateError{
		{
			ErrorCode: "E3002",
			Message:   fmt.Sprintf("conflict: %s with %s = '%s' already exists", table, column, value),
			Field:     column,
			Tag:       "unique",
			Value:     value,
		},
	}
}

// Predefined errors with error codes
var (
	ErrUnauthorized        = CustomError{"E1001", "unauthorized", fiber.StatusUnauthorized}
	ErrForbidden           = CustomError{"E1002", "forbidden", fiber.StatusForbidden}
	ErrNotFound            = CustomError{"E1003", "not found", fiber.StatusNotFound}
	ErrInternal            = CustomError{"E1005", "internal server error", fiber.StatusInternalServerError}
	ErrPaymentRequired     = CustomError{"E1008", "payment required", fiber.StatusPaymentRequired}
	ErrUnprocessableEntity = CustomError{"E1009", "unprocessable entity", fiber.StatusUnprocessableEntity}
	ErrTooManyRequests     = CustomError{"E1010", "too many requests", fiber.StatusTooManyRequests}
	ErrPreconditionFailed  = CustomError{"E1011", "precondition failed", fiber.StatusPreconditionFailed}
)

func ErrBadRequest(err string) CustomError {
	if err == "" {
		err = "bad request"
	}

	return CustomError{
		ErrorCode:  "E1004",
		Message:    err,
		StatusCode: fiber.StatusBadRequest,
	}
}

func ErrInvalidBody(err error) CustomError {
	return CustomError{
		ErrorCode:  "E1006",
		Message:    fmt.Sprintf("invalid request body: %s", err.Error()),
		StatusCode: fiber.StatusBadRequest,
	}
}

// Dynamic error generator for item-specific "not found" errors
func ErrItemNotFound(item string) CustomError {
	return CustomError{
		ErrorCode:  "E2001",
		Message:    fmt.Sprintf("%s not found", item),
		StatusCode: fiber.StatusNotFound,
	}
}

func ErrInvalidItem(item string) CustomError {
	return CustomError{
		ErrorCode:  "E2002",
		Message:    fmt.Sprintf("invalid %s", item),
		StatusCode: fiber.StatusBadRequest,
	}
}

func ErrInvalidParam(param string) CustomError {
	return CustomError{
		ErrorCode:  "E2003",
		Message:    fmt.Sprintf("invalid param %s", param),
		StatusCode: fiber.StatusBadRequest,
	}
}

func ErrDuplicateItem(item string) CustomError {
	return CustomError{
		ErrorCode:  "E2004",
		Message:    fmt.Sprintf("%s already exists", item),
		StatusCode: fiber.StatusConflict,
	}
}
