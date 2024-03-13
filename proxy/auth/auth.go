package auth

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
	"main/proxy/geolocation"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register godoc
// @Summary register users
// @Description Register users using JWT tokens
// @Tags auth
// @Accept  json
// @Produce  json
//
//	@Param			Username and Password				body		string	true	"Username and Password"
//
// @Success 200 {string} string "user is registered"
// @Failure 401 {string} string "error taking a claims"
// @Failure 500 {string} string "user is already exists"
// @Failure 500 {string} string "error hashing password"
// @Router /register [get]
func Register(r chi.Router, users map[string]string) {
	r.Group(func(r chi.Router) {
		r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
			var newUser User
			err := json.NewDecoder(r.Body).Decode(&newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_, ok := users[newUser.Username]
			if ok {
				http.Error(w, "user is already exists", http.StatusInternalServerError)
				return
			}
			hashedPassword, err := bcrypt.
				GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "error hashing password", http.StatusInternalServerError)
				return
			}
			users[newUser.Username] = string(hashedPassword)
			w.Write([]byte("user registered" + newUser.Username + " " + users[newUser.Username] + " "))
		})
	})
}

// Login godoc
// @Summary login users
// @Description login users using JWT tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {string} string "valid JWT token"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "error creating token"
// @Router /login [get]
func LoginandProtectedRoutes(r chi.Router, tokenAuth *jwtauth.JWTAuth, users map[string]string) {
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				http.Error(w, "error", http.StatusUnauthorized)
				return
			}
			password := claims["password"].(string)
			username := claims["username"].(string)
			_, ok := users[username]
			if !ok {
				http.Error(w, "user is not registered", http.StatusUnauthorized)
				return
			}
			_, tokenString, err := tokenAuth.
				Encode(map[string]interface{}{users[username]: password})
			if err != nil {
				http.Error(w, "error creating token", http.StatusInternalServerError)
				return
			} else {
				w.Write([]byte("token " + tokenString))
				w.WriteHeader(http.StatusOK)
				return
			}
		})
		r.Post("/api/address/search", func(w http.ResponseWriter, r *http.Request) {
			geolocation.Search(w, r)
		})
		r.Post("/api/address/geocode", func(w http.ResponseWriter, r *http.Request) {
			geolocation.Geocode(w, r)
		})
	})

}
