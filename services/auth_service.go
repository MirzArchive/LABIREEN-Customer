package services

import (
	"errors"
	"labireen-customer/entities"
	"labireen-customer/pkg/crypto"
	"labireen-customer/pkg/mail"
	"labireen-customer/repositories"
	"os"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterCustomer(customer *entities.CustomerRegister) (mail.EmailData, error)
	LoginCustomer(customer entities.CustomerLogin) (uuid.UUID, error)
	VerifyCustomer(email string) error
	FindByParams(param string, args string) (entities.Customer, error)
	UpdateCustomer(customer entities.CustomerRequest) error
	ResetPassword(pwd entities.CustomerReset, id uuid.UUID) error
}

type authServiceImpl struct {
	rp repositories.CustomerRepository
}

func NewAuthService(rp repositories.CustomerRepository) AuthService {
	return &authServiceImpl{rp}
}

func (svc *authServiceImpl) RegisterCustomer(customer *entities.CustomerRegister) (mail.EmailData, error) {
	if customer.Password != customer.PasswordConfirm {
		return mail.EmailData{}, errors.New("password mismatch")
	}

	hashedPassword, err := crypto.HashValue(customer.Password)
	if err != nil {
		return mail.EmailData{}, errors.New("failed to encrypt given data")
	}

	assignID, err := uuid.NewRandom()
	if err != nil {
		return mail.EmailData{}, errors.New("failed to assign unique uuid")
	}

	user := entities.Customer{
		ID:               assignID,
		Name:             customer.Name,
		Email:            customer.Email,
		Password:         hashedPassword,
		VerificationCode: crypto.Encode(customer.Email),
	}

	err = svc.rp.Create(&user)
	if err != nil {
		return mail.EmailData{}, err
	}

	email := mail.EmailData{
		Email:   []string{user.Email},
		URL:     os.Getenv("EMAIL_CLIENT_ORIGIN") + "/auth/verify/" + user.VerificationCode,
		Name:    user.Name,
		Subject: "Your account verification code",
	}

	return email, nil
}

func (svc *authServiceImpl) LoginCustomer(customer entities.CustomerLogin) (uuid.UUID, error) {
	user, err := svc.rp.GetWhere("email", customer.Email)
	if err != nil {
		return uuid.Nil, errors.New("user not found")
	}

	if !user.Verified {
		return uuid.Nil, errors.New("user has not verified")
	}

	if err := crypto.CheckHash(customer.Password, user.Password); err != nil {
		return uuid.Nil, errors.New("password is not valid or incorrect")
	}

	return user.ID, nil
}

func (svc *authServiceImpl) VerifyCustomer(code string) error {
	user, err := svc.rp.GetWhere("verification_code", code)
	if err != nil {
		return errors.New("user not found")
	}

	user.VerificationCode = ""
	user.Verified = true

	if err := svc.rp.Update(user); err != nil {
		return errors.New("failed to update user data")
	}

	return nil
}

func (svc *authServiceImpl) FindByParams(param string, args string) (entities.Customer, error) {
	user, err := svc.rp.GetWhere(param, args)
	if err != nil {
		return entities.Customer{}, err
	}

	return *user, nil
}

func (svc *authServiceImpl) UpdateCustomer(customer entities.CustomerRequest) error {
	user, err := svc.FindByParams("email", customer.Email)
	if err != nil {
		return err
	}

	user = entities.Customer{
		Name:  customer.Name,
		Email: customer.Email,
		Photo: customer.Photo,
	}

	if err := svc.rp.Update(&user); err != nil {
		return err
	}

	return nil
}

func (svc *authServiceImpl) ResetPassword(pwd entities.CustomerReset, id uuid.UUID) error {
	if pwd.Password != pwd.PasswordConfirm {
		return errors.New("password mismatch")
	}

	user, err := svc.rp.GetById(id)
	if err != nil {
		return err
	}

	hashedPassword, err := crypto.HashValue(pwd.Password)
	if err != nil {
		return errors.New("failed to encrypt given data")
	}

	user.Password = hashedPassword

	if err := svc.rp.Update(user); err != nil {
		return err
	}

	return nil
}
