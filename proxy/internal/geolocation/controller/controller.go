package controller

import (
	"encoding/json"
	"log"
	"main/proxy/internal/geolocation/service"
	"net/http"
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
	geoocodeRequest, err := c.service.GetGeocode(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	geocodeResponse, err := c.service.DadataGeocode(geoocodeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&geocodeResponse)
	if err != nil {
		log.Println("err6")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = c.service.Repository.CacheAddress(*geocodeResponse.Addresses[0])
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

	searchRequest, err := c.service.DadataSearch(r)

	cachedResponse, err := c.service.Repository.GetCache(searchRequest)
	if err == nil {
		err = json.NewEncoder(w).Encode(&cachedResponse)
		if err != nil {
			log.Println("Status 500, dadata.ru is not responding")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	resp, err := c.service.DadataSearchApi(searchRequest)
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
