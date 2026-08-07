package repository

import (
	"database/sql"
	"errors"

	"github.com/vrima25/go-auth-service/interfaces"
	"github.com/vrima25/go-auth-service/model"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) interfaces.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(c *model.Customer) error {
	query := `
	INSERT INTO customers (full_name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.db.QueryRow(query, c.FullName, c.Email, c.PasswordHash).Scan(&c.ID, &c.CreatedAt)
}

func (r *customerRepository) FindByEmail(email string) (*model.Customer, error) {
	query := `
		SELECT id, full_name, email, password_hash, created_at
		FROM customers
		WHERE email = $1
	`

	c := &model.Customer{}
	err := r.db.QueryRow(query, email).Scan(&c.ID, &c.FullName, &c.Email, &c.PasswordHash, &c.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *customerRepository) FindById(ID int) (*model.Customer, error) {
	query := `
		SELECT id, full_name, email, password_hash, created_at
		FROM customers WHERE id = $1
	`
	c := &model.Customer{}
	err := r.db.QueryRow(query, ID).Scan(&c.ID, &c.FullName, &c.Email, &c.PasswordHash, &c.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}
