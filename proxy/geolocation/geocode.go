package geolocation

import (
	"encoding/json"
	"gopkg.in/webdeskltd/dadata.v2"
	"log"
	"net/http"
)

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
//	@Success		200					"Successful operation"
//	@Failure		400				"Bad request"
//	@Failure		401				"Unauthorized"
//	@Failure		404				"Not found"
//	@Failure		500				"Internal server error"
//	@Router			/geocode [post]
func Geocode(w http.ResponseWriter, r *http.Request) {
	var geocodeRequest GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&geocodeRequest)
	if err != nil {
		log.Println("Status")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")

	req := dadata.GeolocateRequest{
		Lat:          geocodeRequest.Lat,
		Lon:          geocodeRequest.Lng,
		Count:        5,
		RadiusMeters: 100,
	}
	addresses, err := api.GeolocateAddress(req)
	if err != nil {
		log.Println("Status 500, dadata.ru is not responding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var geocodeResponse GeocodeResponse
	geocodeResponse.Addresses = []*Address{{Lat: addresses[0].Data.City, Lng: addresses[0].Data.Street + " " + addresses[0].Data.House}}
	err = json.NewEncoder(w).Encode(&geocodeResponse)
	if err != nil {
		log.Println("err5")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
