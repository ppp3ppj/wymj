package middlewaresUsecases

import "github.com/ppp3ppj/wymj/modules/middlewares/middlewaresRepositories"



type IMiddlewaresUsecase interface {
    FindAccessToken(userId, accessToken string) bool
}

type middlewaresUsecase struct {
    middlewaresRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewaresRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
    return &middlewaresUsecase{
        middlewaresRepository: middlewaresRepository,
    }
}

func (u *middlewaresUsecase) FindAccessToken(userId, accessToken string) bool {
    return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}
