package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ppp3ppj/wymj/modules/middlewares/middlewaresHandlers"
	"github.com/ppp3ppj/wymj/modules/middlewares/middlewaresRepositories"
	"github.com/ppp3ppj/wymj/modules/middlewares/middlewaresUsecases"
	"github.com/ppp3ppj/wymj/modules/monitor/monitorHandlers"
	"github.com/ppp3ppj/wymj/modules/users/usersHandlers"
	"github.com/ppp3ppj/wymj/modules/users/usersRepositories"
	"github.com/ppp3ppj/wymj/modules/users/usersUsecases"
)

type IModuleFactory interface {
    MonitorModule()
    UserModule()
}

type moduleFactory struct {
    r fiber.Router // r is router
    s *server // s is server
    mid middlewaresHandlers.IMiddlewaresHandler // mid is middlewaresHandlers
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
    return &moduleFactory{
        r: r,
        s: s,
        mid: mid,
    }
}

func InitMiddleware(s *server) middlewaresHandlers.IMiddlewaresHandler {
    repository := middlewaresRepositories.MiddlewaresRepository(s.db)
    usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
    handler := middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
    return handler
}

func (m *moduleFactory) MonitorModule() {
    handler := monitorHandlers.MonitorHandler(m.s.cfg)

    m.r.Get("/health", handler.HealthCheck)
}

func (m *moduleFactory) UserModule() {
    repository := usersRepositories.UsersRepository(m.s.db)
    usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
    handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

    // Group routes to user = /v1/users/signup
    router := m.r.Group("/users")

    router.Post("/signup", handler.SignUpCustomer)
    router.Post("/signin", handler.SignIn)
    router.Post("/refresh", handler.RefreshPassport)
}
