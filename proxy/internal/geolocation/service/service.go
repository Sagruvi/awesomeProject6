package service

import (
	"gopkg.in/webdeskltd/dadata.v2"
	"log"
	"main/proxy/internal/geolocation/entity"
	"main/proxy/internal/geolocation/repository"
	"strconv"
)

type Service struct {
	Repository repository.Repository
}

func NewService(repository repository.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) CacheSearchHistory(request entity.SearchRequest) error {
	return s.Repository.CacheSearchHistory(request.Query)
}
func (s *Service) CacheAddress(request entity.GeocodeResponse) error {
	return s.Repository.CacheAddress(request)
}
func (s *Service) GetSearchHistory(response entity.SearchResponse) (entity.SearchRequest, error) {
	return s.Repository.GetSearchHistory(response)
}
func (s *Service) GetCache(request entity.SearchRequest) (entity.SearchResponse, error) {
	return s.Repository.GetCache(request.Query)
}

func (s *Service) DadataGeocodeApi(geocodeRequest entity.GeocodeRequest) (string, error) {
	address := entity.Address{
		Lat: strconv.FormatFloat(geocodeRequest.Lat, 'f', -1, 64),
		Lng: strconv.FormatFloat(geocodeRequest.Lng, 'f', -1, 64),
	}
	request := entity.SearchResponse{Addresses: []*entity.Address{&address}}
	cachedResponse, err := s.Repository.GetSearchHistory(request)
	if err == nil {
		return cachedResponse.Query, nil
	}
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")

	req := dadata.GeolocateRequest{
		Lat:          float32(geocodeRequest.Lat),
		Lon:          float32(geocodeRequest.Lng),
		Count:        5,
		RadiusMeters: 100,
	}
	addresses, err := api.GeolocateAddress(req)
	if err != nil {
		log.Println("Status 500, dadata.ru is not responding")
		return "", err
	}

	var geocodeResponse entity.GeocodeResponse
	geocodeResponse.Addresses = []*entity.Address{{Lat: addresses[0].Data.City, Lng: addresses[0].Data.Street + " " + addresses[0].Data.House}}
	res := geocodeResponse.Addresses[0].Lat + " " + geocodeResponse.Addresses[0].Lng
	err = s.CacheAddress(geocodeResponse)
	if err != nil {
		return res, err
	}
	return res, nil

}

func (s *Service) DadataSearchApi(query string) (entity.SearchResponse, error) {
	cachedResponse, err := s.GetCache(entity.SearchRequest{Query: query})
	if err == nil {
		return cachedResponse, nil
	}
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")

	addresses, err := api.CleanAddresses(query)
	if err != nil {
		return entity.SearchResponse{}, err
	}
	log.Println(addresses)
	var searchResponse entity.SearchResponse
	searchResponse.Addresses = []*entity.Address{{Lat: addresses[0].GeoLat, Lng: addresses[0].GeoLon}}
	err = s.CacheSearchHistory(entity.SearchRequest{Query: query})
	return searchResponse, nil
}
