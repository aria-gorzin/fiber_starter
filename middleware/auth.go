package middleware

import (
	"strings"

	"github.com/aria/app/token"
	"github.com/aria/app/util"
	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func Auth(tokenMaker *token.PasetoMaker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			return util.ErrUnauthorized
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return util.ErrUnauthorized
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return util.ErrUnauthorized
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken, token.TokenTypeAuth)
		if err != nil {
			return util.ErrUnauthorized
		}

		ctx.Locals(token.AuthorizationPayloadKey, payload)
		return ctx.Next()
	}
}

func IsAdmin() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := ctx.Locals(token.AuthorizationPayloadKey).(*token.Payload)
		adminRoles := []string{util.SuperuserRole, util.OwnerRole, util.AdminRole, util.OperatorRole, util.DriverRole}
		if util.Includes(adminRoles, payload.Role) {
			return ctx.Next()
		}
		return util.ErrForbidden
	}
}

func IsSuperuser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := ctx.Locals(token.AuthorizationPayloadKey).(*token.Payload)
		if util.SuperuserRole == payload.Role {
			return ctx.Next()
		}
		return util.ErrForbidden
	}
}
