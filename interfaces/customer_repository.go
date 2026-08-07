package interfaces

import "github.com/vrima25/go-auth-service/model"

type CustomerRepository interface {
	Create(customer *model.Customer) error
	FindByEmail(email string) (*model.Customer, error)
	FindById(ID int) (*model.Customer, error)
}
