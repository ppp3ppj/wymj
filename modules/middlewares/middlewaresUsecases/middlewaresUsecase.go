package middlewaresUsecases

import "github.com/ppp3ppj/wymj/modules/middlewares/middlewaresRepositories"



type IMiddlewaresUsecase interface {

}

type middlewaresUsecase struct {
    middlewaresRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewaresRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
    return &middlewaresUsecase{
        middlewaresRepository: middlewaresRepository,
    }
}
