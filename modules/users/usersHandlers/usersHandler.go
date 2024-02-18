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
)

type IUsersHandler interface {
    SignUpCustomer(c *fiber.Ctx) error
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
