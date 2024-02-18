package wymjauth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ppp3ppj/wymj/config"
	"github.com/ppp3ppj/wymj/modules/users"
)

type TokenType string

const (
    Access TokenType = "access"
    Refresh TokenType = "refresh"
    Admin TokenType = "admin"
    ApiKey TokenType = "apiKey"
)

type wymjAuth struct {
    mapClaims *wymjMapClaims // is call payload
    cfg config.IJwtconfig
}

type wymjMapClaims struct {
    Claims *users.UserClaims `json:"claims"`
    jwt.RegisteredClaims 
}

type IWymjAuth interface {
    SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
    return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
    return jwt.NewNumericDate(time.Unix(t, 0))
}

func (w *wymjAuth) SignToken() string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, w.mapClaims)
    ss, _ := token.SignedString(w.cfg.SecretKey())
    return ss
}

func ParseToken(cfg config.IJwtconfig, tokenString string) (*wymjMapClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &wymjMapClaims{
    }, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("signing method is invalid")
        }
        return cfg.SecretKey(), nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenMalformed) {
            return nil, fmt.Errorf("token malformed is invalid")
        } else if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, fmt.Errorf("token had expired")
        } else {
            return nil, fmt.Errorf("parse token failed: %v", err)
        }
    }

    if claims, ok := token.Claims.(*wymjMapClaims); ok {
        return claims, nil
    } else {
        return nil, fmt.Errorf("claims type is invalid")
    }
}

func RepeatToken(cfg config.IJwtconfig, claims *users.UserClaims, exp int64) string {
    obj := &wymjAuth{
        cfg: cfg,
        mapClaims: &wymjMapClaims{
            Claims: claims,
            RegisteredClaims: jwt.RegisteredClaims{
                Issuer: "wymj-api",
                Subject: "refresh-token",
                Audience: []string{"customer", "admin"},
                ExpiresAt: jwtTimeRepeatAdapter(exp),
                NotBefore: jwt.NewNumericDate(time.Now()),
                IssuedAt: jwt.NewNumericDate(time.Now()),
            },
        },
    }
    return obj.SignToken()
}

func NewWymjAuth(tokenType TokenType, cfg config.IJwtconfig, claims *users.UserClaims) (IWymjAuth, error) {
    switch tokenType {
        case Access:
            return newAccessToken(cfg, claims), nil
        case Refresh:
            return newRefreshToken(cfg, claims), nil
        default:
            return nil, fmt.Errorf("unknown token type")
    }
}

func newAccessToken(cfg config.IJwtconfig, claims *users.UserClaims) *wymjAuth {
    return &wymjAuth{
        cfg: cfg,
        mapClaims: &wymjMapClaims{
            Claims: claims,
            RegisteredClaims: jwt.RegisteredClaims{
                Issuer: "wymj-api",
                Subject: "access-token",
                Audience: []string{"customer", "admin"},
                ExpiresAt: jwtTimeDurationCal(cfg.AccessExpireAt()),
                NotBefore: jwt.NewNumericDate(time.Now()),
                IssuedAt: jwt.NewNumericDate(time.Now()),
            },
        },
    }
}

func newRefreshToken(cfg config.IJwtconfig, claims *users.UserClaims) *wymjAuth {
    return &wymjAuth{
        cfg: cfg,
        mapClaims: &wymjMapClaims{
            Claims: claims,
            RegisteredClaims: jwt.RegisteredClaims{
                Issuer: "wymj-api",
                Subject: "refresh-token",
                Audience: []string{"customer", "admin"},
                ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpireAt()),
                NotBefore: jwt.NewNumericDate(time.Now()),
                IssuedAt: jwt.NewNumericDate(time.Now()),
            },
        },
    }
}
