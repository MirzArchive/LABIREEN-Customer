package handlers

import (
	"labireen-customer/entities"
	"labireen-customer/services"
	"os"

	"labireen-customer/pkg/crypto"
	"labireen-customer/pkg/jwtx"
	"labireen-customer/pkg/mail"
	"labireen-customer/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler interface {
	RegisterCustomer(ctx *gin.Context)
	LoginCustomer(ctx *gin.Context)
	VerifyEmail(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
}

type authHandlerImpl struct {
	svc services.AuthService
	ml  mail.EmailSender
}

func NewAuthHandler(svc services.AuthService, ml mail.EmailSender) *authHandlerImpl {
	return &authHandlerImpl{svc, ml}
}

func (hdl *authHandlerImpl) RegisterCustomer(ctx *gin.Context) {
	var request entities.CustomerRegister
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	//Generate Verification Code
	email, err := hdl.svc.RegisterCustomer(&request)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to register user", err.Error())
		return
	}

	hdl.ml.SendEmail(&email)

	response.Success(ctx, http.StatusOK, "User successfuly created, please check your email for email verification", request)
}

func (hdl *authHandlerImpl) LoginCustomer(ctx *gin.Context) {
	var request entities.CustomerLogin
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	id, err := hdl.svc.LoginCustomer(request)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Failed to logged in", err.Error())
		return
	}

	token, err := jwtx.GenerateToken(id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Server error, failed to generate token", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Login Successful", token)
}

func (hdl *authHandlerImpl) VerifyEmail(ctx *gin.Context) {
	code := ctx.Params.ByName("verification-code")

	if err := hdl.svc.VerifyCustomer(code); err != nil {
		response.Error(ctx, http.StatusBadRequest, "User verification failed", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Email verified successfully", nil)
}

func (hdl *authHandlerImpl) ForgotPassword(ctx *gin.Context) {
	var body entities.CustomerRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	user, err := hdl.svc.FindByParams("email", body.Email)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Not Found", err.Error())
	}

	email := mail.EmailData{
		Email:   []string{user.Email},
		URL:     os.Getenv("EMAIL_CLIENT_ORIGIN") + "/auth/resetpassword/" + crypto.Encode(user.ID.String()),
		Name:    user.Name,
		Subject: "Your account reset password link",
	}

	hdl.ml.SendEmail(&email)

	response.Success(ctx, http.StatusOK, "Reset password request successfuly created, please check your email", user)
}

func (hdl *authHandlerImpl) ResetPassword(ctx *gin.Context) {
	token := ctx.Params.ByName("reset-token")

	var body entities.CustomerReset
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	decodedToken, err := crypto.Decode(token)
	if err != nil {
		response.Error(ctx, http.StatusForbidden, "error", err.Error())
		return
	}

	id, err := uuid.Parse(decodedToken)
	if err != nil {
		response.Error(ctx, http.StatusForbidden, "parse error", err.Error())
		return
	}

	if err := hdl.svc.ResetPassword(body, id); err != nil {
		response.Error(ctx, http.StatusForbidden, "reset password failed", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Successfuly change user password", body)
}
