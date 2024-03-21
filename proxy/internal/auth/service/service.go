package service

import (
	"main/proxy/internal/auth/repository"
)

type Service struct {
	Repository repository.Repository
}

func NewService(repository2 repository.Repository) *Service {
	return &Service{
		Repository: repository2,
	}
}

func (s *Service) SaveUser(username, password string) error {
	return s.Repository.SaveUser(username, password)
}

func (s *Service) CheckUser(username, password string) bool {
	return s.Repository.CheckUser(username, password)
}
func (s *Service) CheckPassword(username, password string) bool {
	return s.Repository.CheckPassword(username, password)
}
