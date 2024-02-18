package usersRepositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ppp3ppj/wymj/modules/users"
	"github.com/ppp3ppj/wymj/modules/users/usersPatterns"
)


type IUserRepository interface {
    InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
    FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
}

type userRepository struct {
    db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUserRepository {
    return &userRepository{
        db: db,
    }
}

func (r *userRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
    result := usersPatterns.InsertUser(r.db, req, isAdmin)
    var err error
    if isAdmin {
        result, err = result.Admin()
        if err != nil {
            return nil, err
        }
    } else {
        result, err = result.Customer()
        if err != nil {
            return nil, err
        }
    }
    // Get result from inserting
    user, err := result.Result()
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *userRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
    fmt.Println("email: ", email)
    query := `
    SELECT
        "id",
        "email",
        "password",
        "username",
        "role_id"
    FROM "users"
    WHERE "email" = $1;`

    user := new(users.UserCredentialCheck)
    if err := r.db.Get(user, query, email); err != nil {
        return nil, fmt.Errorf("user not found: %v", err)
    }
    return user, nil
}
