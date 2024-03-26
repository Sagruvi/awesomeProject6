package repository

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db map[string]string
}

func NewRepository() Repository {
	return Repository{
		db: make(map[string]string),
	}
}

func (r *Repository) SaveUser(username, password string) error {
	if r.db[username] != "" {
		return fmt.Errorf("user: %s  already exists", username)
	}
	hashedPassword, err := bcrypt.
		GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	r.db[username] = string(hashedPassword)
	return nil
}
func (r *Repository) CheckUser(username, password string) bool {
	if r.db[username] != "" {
		return true
	}
	return false
}
func (r *Repository) CheckPassword(username, password string) bool {
	hashedPassword := r.db[username]
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
