package middlewaresHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/entities"
	"github.com/ppp3ppj/wymj/modules/middlewares/middlewaresUsecases"
)

type middlewaresHandlerErrCode string

const (
    routerCheckErr middlewaresHandlerErrCode = "middleware-001"
)

type IMiddlewaresHandler interface {
    Cors() fiber.Handler
    RouterCheck() fiber.Handler
}

type middlewaresHandler struct {
    cfg config.IConfig
    middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
    return &middlewaresHandler{
        cfg: cfg,
        middlewaresUsecase: middlewaresUsecase,
    }
}

func (h *middlewaresHandler) Cors() fiber.Handler {
    return cors.New(cors.Config{
        Next: cors.ConfigDefault.Next,
        AllowOrigins: "*",
        AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
        AllowHeaders: "",
        AllowCredentials: false,
        ExposeHeaders: "",
        MaxAge: 0,
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

