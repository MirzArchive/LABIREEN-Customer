package routes

import (
	"labireen-customer/handlers"
	"labireen-customer/middleware"

	"github.com/gin-gonic/gin"
)

type CustomerRoutes struct {
	Router          *gin.Engine
	CustomerHandler handlers.CustomerHandler
}

func (r *CustomerRoutes) Register() {
	customer := r.Router.Group("customer")
	customer.GET("/profile", middleware.ValidateToken(), r.CustomerHandler.GetMe)
}
