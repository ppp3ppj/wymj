package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ppp3ppj/wymj/modules/monitor/monitorHandlers"
)

type IModuleFactory interface {
    MonitorModule()
}

type moduleFactory struct {
    r fiber.Router // r is router
    s *server // s is server
}

func InitModule(r fiber.Router, s *server) IModuleFactory {
    return &moduleFactory{
        r: r,
        s: s,
    }
}

func (m *moduleFactory) MonitorModule() {
    handler := monitorHandlers.MonitorHandler(m.s.cfg)

    m.r.Get("/health", handler.HealthCheck)
}

