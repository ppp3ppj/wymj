package middlewaresHandlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/entities"
	"github.com/ppp3ppj/wymj/modules/middlewares/middlewaresUsecases"
	"github.com/ppp3ppj/wymj/pkg/utils"
	"github.com/ppp3ppj/wymj/pkg/wymjauth"
)

type middlewaresHandlerErrCode string

const (
	routerCheckErr middlewaresHandlerErrCode = "middleware-001"
	jwtAuthErr     middlewaresHandlerErrCode = "middleware-002"
    paramsCheckErr middlewaresHandlerErrCode = "middleware-003"
    authorizeErr   middlewaresHandlerErrCode = "middleware-004"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
    ParamsCheck() fiber.Handler
    Authorize(expectReleId ...int) fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewaresUsecase: middlewaresUsecase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "01/02/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := wymjauth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}
		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
    return func(c *fiber.Ctx) error {
        userId := c.Locals("userId")
        if c.Params("user_id") != userId {
            return entities.NewResponse(c).Error(
                fiber.ErrUnauthorized.Code,
                string(paramsCheckErr),
                "never gonna give you up",
            ).Res()
        }
        return c.Next()
    }
}

func (h *middlewaresHandler) Authorize(expectReleId ...int) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRoleId, ok := c.Locals("userRoleId").(int)
        if !ok {
            return entities.NewResponse(c).Error(
                fiber.ErrUnauthorized.Code,
                string(authorizeErr),
                "user_id is not int type",
            ).Res()
        }
        roles, err := h.middlewaresUsecase.FindRole()
        if err != nil {
            return entities.NewResponse(c).Error(
                fiber.ErrInternalServerError.Code,
                string(authorizeErr),
                err.Error(),
            ).Res()
        }

        sum := 0
        for _, roleId := range expectReleId {
            sum += roleId
        }

        expectedVauleBinary := utils.BinaryConverter(sum, len(roles))
        userVauleBinary := utils.BinaryConverter(userRoleId, len(roles))
        // loop compare bitwise 
        for i := range expectedVauleBinary {
            if userVauleBinary[i]&expectedVauleBinary[i] == 1 {
                return c.Next()
            }
        }
        return entities.NewResponse(c).Error(
            fiber.ErrUnauthorized.Code,
            string(authorizeErr),
            "no permission to access",
        ).Res()
    }
}
