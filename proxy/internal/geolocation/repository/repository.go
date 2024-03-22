package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"log"
	"os"
)

type RepositoryCacher interface {
	CacheSearchHistory(response SearchResponse) error
	CacheAddress(address *Address) error
}

type Repository struct {
	SearchResponse
	SearchRequest
	GeocodeRequest
	GeocodeResponse
	Address
}

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

func NewRepository() *Repository {
	return &Repository{}
}
func (r *Repository) Migrate(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db, err := sql.Open("postgres", string("host="+os.
		Getenv("DB_HOST")+" port="+os.
		Getenv("DB_PORT")+" user="+os.
		Getenv("DB_USER")+" password="+os.
		Getenv("DB_PASSWORD")+" dbname="+os.
		Getenv("DB_NAME")+" sslmode=disable"))
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	defer db.Close()
	if err := goose.Up(db, "migrations/"); err != nil {
		log.Fatalf("Error applying migrations: %s", err)
	}
	return err
}
func Connect(ctx context.Context) (*pgx.Conn, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	conn, err := pgx.Connect(ctx, string("host="+os.
		Getenv("DB_HOST")+" port="+os.
		Getenv("DB_PORT")+" user="+os.
		Getenv("DB_USER")+" password="+os.
		Getenv("DB_PASSWORD")+" dbname="+os.
		Getenv("DB_NAME")+" sslmode=disable"))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
func (r *Repository) CacheSearchHistory(request string) error {
	conn, err := Connect(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "INSERT INTO search_history (query_text) VALUES ($1)", request)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CacheAddress(address Address) error {
	conn, err := Connect(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "INSERT INTO address (lat, lon) VALUES ($1, $2)", address.Lat, address.Lng)
	if err != nil {
		return err
	}

	return nil
}
func (r *Repository) GetSearchHistory(response SearchResponse) (SearchRequest, error) {
	conn, err := Connect(context.Background())
	if err != nil {
		return SearchRequest{}, err
	}
	defer conn.Close(context.Background())

	var searchRequest SearchRequest
	var searchHistoryID int

	err = conn.QueryRow(context.Background(), "SELECT id FROM address WHERE levenshtein(lat, lon, $1, $2) <= LENGTH($1) * 0.3 AND levenshtein(lat, lon, $2, $1) <= LENGTH($2) * 0.3 LIMIT 1;", response.Addresses[0].Lat, response.Addresses[0].Lng).Scan(&searchHistoryID)
	if err != nil {
		return SearchRequest{}, err
	}

	err = conn.QueryRow(context.Background(), "SELECT query_text FROM search_history WHERE id = $1", searchHistoryID).Scan(&searchRequest.Query)
	if err != nil {
		return SearchRequest{}, err
	}

	return searchRequest, nil
}

func (r *Repository) GetCache(request string) (SearchResponse, error) {
	conn, err := Connect(context.Background())
	if err != nil {
		return SearchResponse{}, err
	}
	defer conn.Close(context.Background())

	var searchResponse SearchResponse

	var searchHistoryID int
	err = conn.QueryRow(context.Background(), "SELECT id FROM search_history WHERE levenshtein(query_text, $1) <= LENGTH($1) * 0.3 LIMIT 1", request).Scan(&searchHistoryID)
	if err != nil {
		return SearchResponse{}, err
	}

	rows, err := conn.Query(context.Background(), "SELECT address_text, lat, lon FROM history_search_address JOIN address ON history_search_address.address_id = address.id WHERE history_search_address.search_history_id = $1", searchHistoryID)
	if err != nil {
		return SearchResponse{}, err
	}

	for rows.Next() {
		var address Address
		err = rows.Scan(&address.Lat, &address.Lng)
		if err != nil {
			return SearchResponse{}, err
		}
		searchResponse.Addresses = append(searchResponse.Addresses, &address)
	}

	return searchResponse, nil
}
