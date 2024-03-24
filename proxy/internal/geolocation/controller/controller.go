package controller

import (
	"encoding/json"
	"log"
	"main/proxy/internal/geolocation/entity"
	"main/proxy/internal/geolocation/service"
	"net/http"
)

type Controller struct {
	service service.Service
}
type GeocodeRequest struct {
	entity.GeocodeRequest
}
type SearchRequest struct {
	entity.SearchRequest
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
//	@Param			lat				body		repository.GeocodeRequest	true	"Lat and Lon"
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
	var geocodeRequest GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&geocodeRequest.GeocodeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	geocodeResponse, err := c.service.DadataGeocode(geocodeRequest.GeocodeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&geocodeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = c.service.Repository.CacheAddress(geocodeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&geocodeResponse)
	if err != nil {
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
//	@Param			lat				body		repository.SearchRequest	true	"Address"
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
	var searchRequest SearchRequest
	err := json.NewDecoder(r.Body).Decode(&searchRequest.SearchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	cachedResponse, err := c.service.GetCache(searchRequest.SearchRequest)
	if err == nil {
		err = json.NewEncoder(w).Encode(&cachedResponse)
		if err != nil {
			log.Println("data not found")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	resp, err := c.service.DadataSearchApi(searchRequest.SearchRequest.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Println("Status 500, dadata.ru is not responding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
