package middleware

import (
	"github.com/aria/app/util"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ErrorHandler middleware to handle errors centrally
func Error() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			statusCode := fiber.StatusInternalServerError
			errorMessage := "Internal Server Error"

			// Handle fiber.Error type
			if e, ok := err.(*fiber.Error); ok {
				statusCode = e.Code
				errorMessage = e.Message
				return c.Status(statusCode).JSON(map[string]string{
					"error": errorMessage,
				})
			}

			// Handle CustomError type
			if e, ok := err.(util.CustomError); ok {
				return util.ErrorResponse(c, e)
			}

			// Handle validation errors from validator
			if validationErrors, ok := err.(util.ValidateErrors); ok {
				return util.ValidationErrorResponse(c, validationErrors)
			}

			log.Error().Err(err).Msg("An unexpected error occurred")
			// Handle unknown errors
			return c.Status(statusCode).JSON(map[string]string{
				"error": errorMessage,
			})
		}
		return nil
	}
}
