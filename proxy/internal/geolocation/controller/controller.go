package controller

import (
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
	c.service.DadataGeocode(w, r)
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
	c.service.DadataSearch(w, r)
}
