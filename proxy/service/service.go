package service

import (
	"github.com/go-chi/jwtauth"
	"main/proxy/repository"
)

type Service struct {
	Repository repository.Repository
	TokenAuth  *jwtauth.JWTAuth
}

func NewService(secret string) *Service {
	return &Service{
		Repository: repository.NewRepository(),
		TokenAuth:  jwtauth.New("HS256", []byte(secret), nil),
	}
}
