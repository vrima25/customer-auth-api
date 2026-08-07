package service

import (
	"errors"

	"github.com/vrima25/go-auth-service/interfaces"
	"github.com/vrima25/go-auth-service/model"
	"github.com/vrima25/go-auth-service/util"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid email or password")
)

type authService struct {
	customerRepo interfaces.CustomerRepository
	jwtSecret    string
}

func NewAuthService(customerRepo interfaces.CustomerRepository, jwtSecret string) interfaces.AuthService {
	return &authService{customerRepo: customerRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(email, password, full_name string) (*model.Customer, error) {
	existing, err := s.customerRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	customer := &model.Customer{
		Email:        email,
		PasswordHash: string(hashed),
		FullName:     full_name,
	}

	if err := s.customerRepo.Create(customer); err != nil {
		return nil, err
	}

	return customer, nil

}

func (s *authService) Login(email, password string) (string, *model.Customer, error) {
	customer, err := s.customerRepo.FindByEmail(email)

	if err != nil {
		return "", nil, err
	}

	if customer == nil {
		return "", nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(password))

	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := util.GenerateToken(customer.Email, s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, customer, nil
}
