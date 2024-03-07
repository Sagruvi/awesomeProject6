package geolocation

import (
	"encoding/json"
	"gopkg.in/webdeskltd/dadata.v2"
	"log"
	"net/http"
)

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
