package appinfohandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/appinfo/appinfoUsecases"
	"github.com/ppp3ppj/wymj/modules/entities"
	"github.com/ppp3ppj/wymj/pkg/wymjauth"
)

type appinfoHandlersErrorCode string

const (
    generateApiKeyError appinfoHandlersErrorCode = "appinfo-001"
)

type IAppinfoHandler interface {
    GenerateApiKey(c *fiber.Ctx) error
}

type appinfoHandler struct {
    cfg config.IConfig
    appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
    return &appinfoHandler{
        cfg: cfg,
        appinfoUsecase: appinfoUsecase,
    }
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
    apiKey, err := wymjauth.NewWymjAuth(
       wymjauth.ApiKey,
       h.cfg.Jwt(),
       nil,
    )

    if err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrInternalServerError.Code,
            string(generateApiKeyError),
            err.Error(),
        ).Res()
    }

    return entities.NewResponse(c).Success(
        fiber.StatusOK,
        &struct {
            Key string `json:"key"`
        } {
            Key: apiKey.SignToken(),
        },
    ).Res()
}
