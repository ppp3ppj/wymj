package main

import (
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
    db := databases.DbConnect(cfg.Db())
    servers.NewServer(cfg, db).Start()
}
