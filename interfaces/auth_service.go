package interfaces

import "github.com/vrima25/go-auth-service/model"

type AuthService interface {
	Register(email, password, fullName string) (*model.Customer, error)
	Login(email, password string) (token string, customer *model.Customer, err error)
}
