package appinfohandlers

import (
	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/appinfo/appinfoUsecases"
)

type IAppinfoHandler interface {

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
