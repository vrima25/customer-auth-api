package interfaces

import "ticket-triage-api/model"

type CustomerRepository interface {
	Create(customer *model.Customer) error
	FindByEmail(email string) (*model.Customer, error)
	FindById(ID int) (*model.Customer, error)
}