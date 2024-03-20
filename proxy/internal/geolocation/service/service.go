package service

import "main/proxy/internal/geolocation/repository"

type Service struct {
	Repository repository.Repository
}

func NewService(repository repository.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) CacheSearchHistory(request repository.SearchRequest) error {
	return s.Repository.CacheSearchHistory(request)
}
func (s *Service) CacheAddress(address repository.Address) error {
	return s.Repository.CacheAddress(address)
}
func (s *Service) GetSearchHistory(response repository.SearchResponse) (repository.SearchRequest, error) {
	return s.Repository.GetSearchHistory(response)
}
func (s *Service) GetCache(request repository.SearchRequest) (repository.SearchResponse, error) {
	return s.Repository.GetCache(request)
}
