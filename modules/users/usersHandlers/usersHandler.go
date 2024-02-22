package usersHandlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/entities"
	"github.com/ppp3ppj/wymj/modules/users"
	"github.com/ppp3ppj/wymj/modules/users/usersUsecases"
	"github.com/ppp3ppj/wymj/pkg/wymjauth"
)

type userHandlerErrCode string 
const (
    signupCustomerErr userHandlerErrCode = "users-001"
    signInErr userHandlerErrCode = "users-002"
    refreshPassportErr userHandlerErrCode = "users-003"
    signOutErr userHandlerErrCode = "users-004"
    singupAdminErr userHandlerErrCode = "users-005"
    generateAdminTokenErr userHandlerErrCode = "users-006"
    getUserProfileErr userHandlerErrCode = "users-007"
)

type IUsersHandler interface {
    SignUpCustomer(c *fiber.Ctx) error
    SignIn(c *fiber.Ctx) error
    RefreshPassport(c *fiber.Ctx) error
    SignOut(c *fiber.Ctx) error
    SignUpAdmin(c *fiber.Ctx) error
    GenerateAdminToken(c *fiber.Ctx) error
    GetUserProfile(c *fiber.Ctx) error
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

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
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

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
    adminToken, err := wymjauth.NewWymjAuth(
        wymjauth.Admin,
        h.cfg.Jwt(),
        nil,
    )
    if err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrInternalServerError.Code,
            string(generateAdminTokenErr),
            err.Error(),
        ).Res()
    }

    return entities.NewResponse(c).Success(fiber.StatusOK, &struct{
        Token string `json:"token"`
    }{
        Token: adminToken.SignToken(),
    }).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
    userId := strings.Trim(c.Params("user_id"), " ")

    result, err := h.usersUsecase.GetUserProfile(userId)
    if err != nil {
        switch err.Error() {
        case "get user faild: sql: no rows in result set":
            return entities.NewResponse(c).Error(
                fiber.ErrNotFound.Code,
                string(getUserProfileErr),
                err.Error(),
            ).Res()
        default:
            return entities.NewResponse(c).Error(
                fiber.ErrInternalServerError.Code,
                string(getUserProfileErr),
                err.Error(),
            ).Res()
        }
    }
    return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
