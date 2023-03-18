package repositories

import (
	"labireen-customer/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *entities.Customer) error
	GetById(id uuid.UUID) (*entities.Customer, error)
	GetWhere(param string, email string) (*entities.Customer, error)
	Update(customer *entities.Customer) error
	Delete(customer *entities.Customer) error
}

type customerRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepositoryImpl{db}
}

func (rp *customerRepositoryImpl) Create(customer *entities.Customer) error {
	return rp.db.Create(&customer).Error
}

func (rp *customerRepositoryImpl) GetById(id uuid.UUID) (*entities.Customer, error) {
	var customer entities.Customer
	if err := rp.db.First(&customer, id).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

func (rp *customerRepositoryImpl) GetWhere(param string, args string) (*entities.Customer, error) {
	var customer entities.Customer
	if err := rp.db.Where(param+" = ?", args).First(&customer).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

func (rp *customerRepositoryImpl) Update(customer *entities.Customer) error {
	return rp.db.Save(customer).Error
}

func (rp *customerRepositoryImpl) Delete(customer *entities.Customer) error {
	return rp.db.Delete(customer).Error
}
