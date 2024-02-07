package main

import (
	"fmt"
	"os"

	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/servers"
	"github.com/ppp3ppj/wymj/pkg/databases"
)

func envPath() string {
    if len(os.Args) == 1 {
        return ".env"
    } else {
        return os.Args[1]
    }
}

func main() {
    cfg := config.LoadConfig(envPath())
    fmt.Println(cfg.App())
    fmt.Println(cfg.Db())
    fmt.Println(cfg.Jwt())
    db := databases.DbConnect(cfg.Db())
    fmt.Println(db)

    servers.NewServer(cfg, db).Start()
}
