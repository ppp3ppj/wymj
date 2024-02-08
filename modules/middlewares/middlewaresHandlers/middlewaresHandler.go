package middlewaresHandlers

import "github.com/ppp3ppj/wymj/modules/middlewares/middlewaresUsecases"


type IMiddlewaresHandler interface {

}

type middlewaresHandler struct {
    middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(middlewaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
    return &middlewaresHandler{
        middlewaresUsecase: middlewaresUsecase,
    }
}
