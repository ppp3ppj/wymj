package usersHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/entities"
	"github.com/ppp3ppj/wymj/modules/users"
	"github.com/ppp3ppj/wymj/modules/users/usersUsecases"
)

type userHandlerErrCode string 
const (
    signupCustomerErr userHandlerErrCode = "users-001"
    signInErr userHandlerErrCode = "users-002"
    refreshPassportErr userHandlerErrCode = "users-003"
    signOutErr userHandlerErrCode = "users-004"
)

type IUsersHandler interface {
    SignUpCustomer(c *fiber.Ctx) error
    SignIn(c *fiber.Ctx) error
    RefreshPassport(c *fiber.Ctx) error
    SignOut(c *fiber.Ctx) error
}

type usersHandler struct {
    cfg config.IConfig
    usersUsecase usersUsecases.IUserUsecase
}

func UsersHandler(cfg config.IConfig, userUsecase usersUsecases.IUserUsecase) IUsersHandler {
    return &usersHandler{
        cfg: cfg,
        usersUsecase: userUsecase,
    }
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
    // Request body parsing
    req := new(users.UserRegisterReq)
    if err := c.BodyParser(req); err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(signupCustomerErr),
            err.Error(),
        ).Res()
    }
    // Email validation
    if !req.IsEmail() {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(signupCustomerErr),
            "email pattern is invalid",
        ).Res()
    }
    // Insert users
    result, err := h.usersUsecase.InsertCustomer(req)
    if err != nil {
        switch err.Error() {
            case "username has been used": 
                return entities.NewResponse(c).Error(
                    fiber.ErrBadRequest.Code,
                    string(signupCustomerErr),
                    err.Error(),
                ).Res()
            case "email has been used":
                return entities.NewResponse(c).Error(
                    fiber.ErrBadRequest.Code,
                    string(signupCustomerErr),
                    err.Error(),
                ).Res()
            default:
                return entities.NewResponse(c).Error(
                    fiber.ErrInternalServerError.Code,
                    string(signupCustomerErr),
                    err.Error(),
                ).Res()
        }
    }
    return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
    req := new(users.UserCredential)
    if err := c.BodyParser(req); err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(signInErr),
            err.Error(),
        ).Res()
    }

    passport, err := h.usersUsecase.GetPassport(req)
    if err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrUnauthorized.Code,
            string(signInErr),
            err.Error(),
        ).Res()
    }
    return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
    req := new(users.UserRefreshCredential)
    if err := c.BodyParser(req); err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(refreshPassportErr),
            err.Error(),
        ).Res()
    }

    passport, err := h.usersUsecase.RefreshPassport(req)
    if err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrUnauthorized.Code,
            string(refreshPassportErr),
            err.Error(),
        ).Res()
    }
    return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
    req := new(users.UserRemoveCredential)

    if err := c.BodyParser(req); err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(signOutErr),
            err.Error(),
        ).Res()
    }
    
    if err := h.usersUsecase.DeleteOauth(req.OauthId); err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(signOutErr),
            err.Error(),
        ).Res()
    }

    return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
