package entity

type SearchRequest struct {
	Query string `json:"query"`
}
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}
type GeocodeRequest struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lon"`
}
type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}
type Address struct {
	Lat string `json:"lat"`
	Lng string `json:"lon"`
}
