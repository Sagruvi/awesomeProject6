package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"main/proxy/auth"
	_ "main/proxy/docs" // docs is generated by Swag CLI, you have to import it.
	"net/http"
	"net/http/httputil"
	"net/url"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = make(map[string]string)
var tokenAuth *jwtauth.JWTAuth

//	@title			Dadata API Proxy
//	@version		1.0
//	@description	This is a sample server geolocation service.

// @host	localhost:8080/
// @BasePath	/api/address
func main() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	hugoURL, err := url.Parse("http://hugo:1313")
	if err != nil {
		panic(err)
	}
	hugoProxy := httputil.NewSingleHostReverseProxy(hugoURL)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		))

		r.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from API"))
		})
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			hugoProxy.ServeHTTP(w, r)
		})
	})
	auth.Register(r, tokenAuth, users)
	auth.Login(r, tokenAuth, users)
	auth.ProtectedRoutes(r, tokenAuth)
	http.ListenAndServe(":8080", r)
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

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://"+rp.host+":"+rp.port+r.RequestURI, http.StatusFound)
}

// @contact.name	API Support
// @contact.url	https://github.com/go-chi/chi/issues
// @contact.email	6z6o8@example.com
