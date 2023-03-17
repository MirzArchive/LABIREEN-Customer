package services

import (
	"labireen-customer/entities"
	"labireen-customer/repositories"

	"github.com/google/uuid"
)

type CustomerService interface {
	UpdateCustomer(customer entities.CustomerRequest) error
	GetCustomer(id uuid.UUID) (entities.CustomerRequest, error)
}

type customerServiceImpl struct {
	repo repositories.AuthRepository
}

func NewCustomerService(repo repositories.AuthRepository) CustomerService {
	return &customerServiceImpl{repo}
}

func (csr *customerServiceImpl) UpdateCustomer(customer entities.CustomerRequest) error {
	return nil
}

func (csr *customerServiceImpl) GetCustomer(id uuid.UUID) (entities.CustomerRequest, error) {
	user, err := csr.repo.GetById(id)
	if err != nil {
		return entities.CustomerRequest{}, err
	}

	userResp := entities.CustomerRequest{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userResp, nil
}
