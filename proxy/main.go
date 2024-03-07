package main

import (
	"main/proxy/geolocation"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {

	// Создаем реверсивный прокси для сервера Hugo
	hugoURL, err := url.Parse("http://hugo:1313")
	if err != nil {
		panic(err)
	}
	hugoProxy := httputil.NewSingleHostReverseProxy(hugoURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/address/search") {
			geolocation.Search(w, r)
		}
		if strings.HasPrefix(r.URL.Path, "/api/address/geocode") {
			geolocation.Geocode(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Write([]byte("Hello from API"))
		} else {
			hugoProxy.ServeHTTP(w, r)
		}
	})
	//go buildmermaid.WorkerBinary()
	//go buildmermaid.WorkerGraph()
	// Запускаем сервер на порту 8080
	http.ListenAndServe(":8080", nil)
}

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
