package controller

import (
	"main/proxy/internal/service"
	"net/http"
)

type Controller struct {
	Service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		Service: service,
	}
}

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	c.Service.Register(w, r)
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	c.Service.Login(w, r)
}
