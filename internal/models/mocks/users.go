package mocks

import (
	"github.com/mixnblend/snippetbox/internal/models"
)

type UserModel struct{}

type UserCredentials struct {
	UserName string
	Password string
}

const DuplicateEmail = "dupe@example.com"

var ValidUserCredentials = &UserCredentials{
	UserName: "alice@example.com",
	Password: "pa$$word",
}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case DuplicateEmail:
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == ValidUserCredentials.UserName && password == ValidUserCredentials.Password {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
