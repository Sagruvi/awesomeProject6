package geolocation

import (
	"encoding/json"
	"gopkg.in/webdeskltd/dadata.v2"
	"log"
	"net/http"
)

type SearchRequest struct {
	Query string `json:"query"`
}
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}
type GeocodeRequest struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lon"`
}
type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
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
func Search(w http.ResponseWriter, r *http.Request) {
	var searchRequest SearchRequest
	err := json.NewDecoder(r.Body).Decode(&searchRequest)
	if err != nil {
		log.Println("err502.1")
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	var searchResponse SearchResponse
	searchResponse.Addresses = []*Address{{Lat: addresses[0].GeoLat, Lng: addresses[0].GeoLon}}
	err = json.NewEncoder(w).Encode(&searchResponse)
	if err != nil {
		log.Println("Status 500, dadata.ru is not responding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Address struct {
	Lat string `json:"lat"`
	Lng string `json:"lon"`
}
