package usersUsecases

import (
	"fmt"

	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/users"
	"github.com/ppp3ppj/wymj/modules/users/usersRepositories"
	"github.com/ppp3ppj/wymj/pkg/wymjauth"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
    InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
    GetPassport(req *users.UserCredential) (*users.UserPassport, error)
}

type userUsecase struct {
    cfg config.IConfig
    userRepository usersRepositories.IUserRepository
}

func UsersUsecase(cfg config.IConfig, userRepository usersRepositories.IUserRepository) IUserUsecase {
    return &userUsecase{
        cfg: cfg,
        userRepository: userRepository,
    }
}

func (u *userUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
    // Hashing password
    if err := req.BcryptHashing(); err != nil {
        return nil, err
    }
    // Insert user
    result, err := u.userRepository.InsertUser(req, false)
    if err != nil {
        return nil, err
    }
    return result, nil
}

func (u *userUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
    // Find user
    user, err := u.userRepository.FindOneUserByEmail(req.Email)
    if err != nil {
        return nil, err
    }

    // Compare password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, fmt.Errorf("password is invalid")
    }

    accessToken, err := wymjauth.NewWymjAuth(wymjauth.Access, u.cfg.Jwt(), &users.UserClaims{
        Id: user.Id,
        RoleId: user.RoleId,
    })

    refreshToken, err := wymjauth.NewWymjAuth(wymjauth.Refresh, u.cfg.Jwt(), &users.UserClaims{
        Id: user.Id,
        RoleId: user.RoleId,
    })

    // Set user passport
    passport := &users.UserPassport{
        User: &users.User{
            Id: user.Id,
            Email: user.Email,
            Username: user.Username,
            RoleId: user.RoleId,
        },
        Token: &users.UserToken{
            AccessToken: accessToken.SignToken(),
            RefreshToken: refreshToken.SignToken(),
        },
    }

    if err := u.userRepository.InsertOauth(passport); err != nil {
        return nil, err
    }
    return passport, nil
}
