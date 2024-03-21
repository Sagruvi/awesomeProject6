package service

import (
	"encoding/json"
	"gopkg.in/webdeskltd/dadata.v2"
	"log"
	"main/proxy/internal/geolocation/repository"
	"net/http"
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
func (s *Service) CacheAddress(address repository.Address) error {
	return s.Repository.CacheAddress(address)
}
func (s *Service) GetSearchHistory(response repository.SearchResponse) (repository.SearchRequest, error) {
	return s.Repository.GetSearchHistory(response)
}
func (s *Service) GetCache(request repository.SearchRequest) (repository.SearchResponse, error) {
	return s.Repository.GetCache(request.Query)
}
func (s *Service) DadataGeocode(w http.ResponseWriter, r *http.Request) {
	var geocodeRequest repository.GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&geocodeRequest)
	if err != nil {
		log.Println("Status")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	address := repository.Address{
		Lat: strconv.FormatFloat(geocodeRequest.Lat, 'f', -1, 64),
		Lng: strconv.FormatFloat(geocodeRequest.Lng, 'f', -1, 64),
	}
	request := repository.SearchResponse{Addresses: []*repository.Address{&address}}
	cachedResponse, err := s.Repository.GetSearchHistory(request)
	if err == nil {
		err = json.NewEncoder(w).Encode(&cachedResponse)
		if err != nil {
			log.Println("err5")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var geocodeResponse repository.GeocodeResponse
	geocodeResponse.Addresses = []*repository.Address{{Lat: addresses[0].Data.City, Lng: addresses[0].Data.Street + " " + addresses[0].Data.House}}

	err = json.NewEncoder(w).Encode(&geocodeResponse)
	if err != nil {
		log.Println("err6")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.Repository.CacheAddress(*geocodeResponse.Addresses[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&geocodeResponse)
	if err != nil {
		log.Println("err5")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func (s *Service) DadataSearch(w http.ResponseWriter, r *http.Request) {
	var searchRequest repository.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&searchRequest)
	if err != nil {
		log.Println("err502.1")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cachedResponse, err := s.Repository.GetCache(searchRequest.Query)
	if err == nil {
		err = json.NewEncoder(w).Encode(&cachedResponse)
		if err != nil {
			log.Println("Status 500, dadata.ru is not responding")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")

	log.Println(searchRequest.Query)

	addresses, err := api.CleanAddresses(searchRequest.Query)
	if err != nil {
		log.Println("err502.2")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(addresses)
	var searchResponse repository.SearchResponse
	searchResponse.Addresses = []*repository.Address{{Lat: addresses[0].GeoLat, Lng: addresses[0].GeoLon}}
	err = json.NewEncoder(w).Encode(&searchResponse)
	if err != nil {
		log.Println("Status 500, dadata.ru is not responding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
