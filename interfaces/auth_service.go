package interfaces

import "ticket-triage-api/model"

type AuthService interface {
	Register(email, password, fullName string) (*model.Customer, error)
	Login(email, password string) (token string, customer *model.Customer, err error)
}