package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/ppp3ppj/wymj/config"
)

type IServer interface {
    Start()
}

type server struct {
    app *fiber.App
    cfg config.IConfig
    db *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
    app := fiber.New(fiber.Config {
        AppName: cfg.App().Name(),
        BodyLimit: cfg.App().BodyLimit(),
        ReadTimeout: cfg.App().ReadTimeout(),
        WriteTimeout: cfg.App().WriteTimeout(),
        JSONEncoder: json.Marshal,
        JSONDecoder: json.Unmarshal,
    })
    return &server{
        cfg: cfg,
        db: db,
        app: app,
    }
}

func (s *server) Start() {
    // Middlewares
    middlewares := InitMiddleware(s)
    s.app.Use(middlewares.Logger())
    s.app.Use(middlewares.Cors())
    // Modules
    // http://localhost:3000/v1
    v1 := s.app.Group("/v1")
    modules := InitModule(v1, s, middlewares)
    modules.MonitorModule()

    s.app.Use(middlewares.RouterCheck())
    // Graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func() {
        _ = <-c
        log.Panicln("Shutting down server...")
        _ = s.app.Shutdown()
    }()

    // Listen to host:port
    log.Printf("Server is running on %v", s.cfg.App().Url())
    s.app.Listen(s.cfg.App().Url())
}
