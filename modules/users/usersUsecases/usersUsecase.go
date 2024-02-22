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
    InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
    GetPassport(req *users.UserCredential) (*users.UserPassport, error)
    RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
    DeleteOauth(oauthId string) error
    GetUserProfile(userId string) (*users.User, error)
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

func (u *userUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
    // Hashing password
    if err := req.BcryptHashing(); err != nil {
        return nil, err
    }
    // Insert user
    result, err := u.userRepository.InsertUser(req, true)
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

func (u *userUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
    // Parse Token
    claims, err := wymjauth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
    if err != nil {
        return nil, err
    }

    // Find oauth
    oauth, err := u.userRepository.FindOneOauth(req.RefreshToken)
    if err != nil {
        return nil, err
    }

    // Find user profile
    profile, err := u.userRepository.GetProfile(oauth.UserId)

    newClaims := &users.UserClaims{
        Id: profile.Id,
        RoleId: profile.RoleId,
    }

    accessToken, err := wymjauth.NewWymjAuth(
        wymjauth.Access, 
        u.cfg.Jwt(), 
        newClaims,
    )
    
    refreshToken := wymjauth.RepeatToken(
        u.cfg.Jwt(), 
        newClaims,
        claims.ExpiresAt.Unix(),
    )

    passport := &users.UserPassport{
        User: profile,
        Token: &users.UserToken{
            Id: oauth.Id,
            AccessToken: accessToken.SignToken(),
            RefreshToken: refreshToken,
        },
    }

    if err := u.userRepository.UpdateOauth(passport.Token); err != nil {
        return nil, err
    }
    return passport, nil
}

func (u *userUsecase) DeleteOauth(oauthId string) error {
    if err := u.userRepository.DeleteOauth(oauthId); err != nil {
        return err
    }
    return nil
}

func (u *userUsecase) GetUserProfile(userId string) (*users.User, error) {
    profile, err := u.userRepository.GetProfile(userId)
    if err != nil {
        return nil, err
    }
    return profile, nil
}
