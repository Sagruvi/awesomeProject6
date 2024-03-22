package service

import (
	"encoding/json"
	"gopkg.in/webdeskltd/dadata.v2"
	"io"
	"log"
	"main/proxy/internal/geolocation/repository"
	"strconv"
)

type Service struct {
	Repository repository.Repository
}

func NewService(repository repository.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) CacheSearchHistory(request repository.SearchRequest) error {
	return s.Repository.CacheSearchHistory(request.Query)
}
func (s *Service) CacheAddress(request repository.GeocodeResponse) error {
	return s.Repository.CacheAddress(request)
}
func (s *Service) GetSearchHistory(response repository.SearchResponse) (repository.SearchRequest, error) {
	return s.Repository.GetSearchHistory(response)
}
func (s *Service) GetCache(request repository.SearchRequest) (repository.SearchResponse, error) {
	return s.Repository.GetCache(request.Query)
}
func (s *Service) GetGeocode(r io.ReadCloser) (repository.GeocodeRequest, error) {
	var geocodeRequest repository.GeocodeRequest
	err := json.NewDecoder(r).Decode(&geocodeRequest)
	if err != nil {
		log.Println("Status")
		return geocodeRequest, err
	}
	return geocodeRequest, nil
}
func (s *Service) DadataGeocodeApi(geocodeRequest repository.GeocodeRequest) (repository.SearchRequest, error) {
	address := repository.Address{
		Lat: strconv.FormatFloat(geocodeRequest.Lat, 'f', -1, 64),
		Lng: strconv.FormatFloat(geocodeRequest.Lng, 'f', -1, 64),
	}
	request := repository.SearchResponse{Addresses: []*repository.Address{&address}}
	cachedResponse, err := s.Repository.GetSearchHistory(request)
	return cachedResponse, err

}
func (s *Service) DadataGeocode(geocodeRequest repository.GeocodeRequest) (repository.GeocodeResponse, error) {

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
		return repository.GeocodeResponse{}, err
	}

	var geocodeResponse repository.GeocodeResponse
	geocodeResponse.Addresses = []*repository.Address{{Lat: addresses[0].Data.City, Lng: addresses[0].Data.Street + " " + addresses[0].Data.House}}
	return geocodeResponse, nil

}
func (s *Service) DadataSearch(r io.ReadCloser) (string, error) {
	var searchRequest repository.SearchRequest
	err := json.NewDecoder(r).Decode(&searchRequest)
	if err != nil {
		return "", err
	}
	return searchRequest.Query, nil
}
func (s *Service) DadataSearchApi(query string) (repository.SearchResponse, error) {
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")

	addresses, err := api.CleanAddresses(query)
	if err != nil {
		return repository.SearchResponse{}, err
	}
	log.Println(addresses)
	var searchResponse repository.SearchResponse
	searchResponse.Addresses = []*repository.Address{{Lat: addresses[0].GeoLat, Lng: addresses[0].GeoLon}}
	return searchResponse, nil
}
