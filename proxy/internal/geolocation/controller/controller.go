package controller

import (
	"encoding/json"
	"gopkg.in/webdeskltd/dadata.v2"
	"log"
	"main/proxy/internal/geolocation/repository"
	"main/proxy/internal/geolocation/service"
	"net/http"
	"strconv"
)

type Controller struct {
	service service.Service
}

func NewController(service2 service.Service) Controller {
	return Controller{service: service2}
}

// SearchAddress godoc
//
//	@Summary		Search for address suggestions
//	@Description	Search for address suggestions by latitude and longitude
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Param			lat				body		GeocodeRequest	true	"Lat and Lon"
//	@Param			Authorization	header		string			true	"Authorization token"
//	@Param			X-Secret		header		string			true	"API Private token"
//
// @Param        Authorization header string true "Bearer token"
//
//	@Success		200					"Successful operation"
//	@Failure		400				"Bad request"
//	@Failure		401				"Unauthorized"
//	@Failure		404				"Not found"
//	@Failure		500				"Internal server error"
//	@Router			/geocode [post]
func (c *Controller) Geocode(w http.ResponseWriter, r *http.Request) {
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
	cachedResponse, err := c.service.Repository.GetSearchHistory(request)
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
	err = c.service.CacheAddress(*geocodeResponse.Addresses[0])
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

// SearchAddress godoc
//
//	@Summary		Search for address
//	@Description	Search for latitude and longitude by address
//	@Tags			addresses
//	@Accept			json
//	@Produce		json
//	@Param			lat				body		SearchRequest	true	"Address"
//	@Param			Authorization	header		string			true	"Authorization token"
//	@Param			X-Secret		header		string			true	"API Private token"
//
// @Param        Authorization header string true "Bearer token"
//
//	@Success		200					"Successful operation"
//	@Failure		400				"Bad request"
//	@Failure		401				"Unauthorized"
//	@Failure		404				"Not found"
//	@Failure		500				"Internal server error"
//	@Router			/search [post]
func (c *Controller) Search(w http.ResponseWriter, r *http.Request) {
	var searchRequest repository.SearchRequest
	err := json.NewDecoder(r.Body).Decode(&searchRequest)
	if err != nil {
		log.Println("err502.1")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cachedResponse, err := c.service.Repository.GetCache(searchRequest.Query)
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
