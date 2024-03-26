package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"log"
	"main/proxy/internal/geolocation/entity"
	"os"
)

type RepositoryCacher interface {
	CacheSearchHistory(response entity.SearchResponse) error
	CacheAddress(address *entity.Address) error
}

type Repository struct {
	client *redis.Client
}

func NewRepository() *Repository {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &Repository{client: client}
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
	redisErr := r.client.Set(context.Background(), "search_history"+request, request, 0)
	if redisErr.Err() != nil {
		return redisErr.Err()
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO search_history (query_text) VALUES ($1)", request)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CacheAddress(geocodeResponse entity.GeocodeResponse) error {
	conn, err := Connect(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	redisErr := r.client.Set(context.Background(), "address"+geocodeResponse.Addresses[0].Lat+" "+geocodeResponse.Addresses[0].Lng, geocodeResponse, 0)
	if redisErr.Err() != nil {
		return redisErr.Err()
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO address (lat, lon) VALUES ($1, $2)", geocodeResponse.Addresses[0].Lat, geocodeResponse.Addresses[0].Lng)
	if err != nil {
		return err
	}

	return nil
}
func (r *Repository) GetSearchHistory(response entity.SearchResponse) (entity.SearchRequest, error) {
	conn, err := Connect(context.Background())
	if err != nil {
		return entity.SearchRequest{}, err
	}
	defer conn.Close(context.Background())

	var searchRequest entity.SearchRequest
	var searchHistoryID int
	res := r.client.Get(context.Background(), "search_history"+response.Addresses[0].Lat+" "+response.Addresses[0].Lng)
	if res.Err() != nil {
		return entity.SearchRequest{}, res.Err()
	}
	if res.Val() != "" {
		err = json.Unmarshal([]byte(res.Val()), &searchRequest)
		if err != nil {
			return entity.SearchRequest{}, err
		}
	}
	err = conn.QueryRow(context.Background(), "SELECT id FROM address WHERE levenshtein(lat, lon, $1, $2) <= LENGTH($1) * 0.3 AND levenshtein(lat, lon, $2, $1) <= LENGTH($2) * 0.3 LIMIT 1;", response.Addresses[0].Lat, response.Addresses[0].Lng).Scan(&searchHistoryID)
	if err != nil {
		return entity.SearchRequest{}, err
	}

	err = conn.QueryRow(context.Background(), "SELECT query_text FROM search_history WHERE id = $1", searchHistoryID).Scan(&searchRequest.Query)
	if err != nil {
		return entity.SearchRequest{}, err
	}

	return searchRequest, nil
}

func (r *Repository) GetCache(request string) (entity.SearchResponse, error) {
	conn, err := Connect(context.Background())
	if err != nil {
		return entity.SearchResponse{}, err
	}
	defer conn.Close(context.Background())

	var searchResponse entity.SearchResponse

	var searchHistoryID int
	res := r.client.Get(context.Background(), "address"+request)
	if res.Err() != nil {
		return entity.SearchResponse{}, res.Err()
	}
	if res.Val() != "" {
		err = json.Unmarshal([]byte(res.Val()), &searchResponse)
		if err != nil {
			return entity.SearchResponse{}, err
		}
	}
	err = conn.QueryRow(context.Background(), "SELECT id FROM search_history WHERE levenshtein(query_text, $1) <= LENGTH($1) * 0.3 LIMIT 1", request).Scan(&searchHistoryID)
	if err != nil {
		return entity.SearchResponse{}, err
	}

	rows, err := conn.Query(context.Background(), "SELECT address_text, lat, lon FROM history_search_address JOIN address ON history_search_address.address_id = address.id WHERE history_search_address.search_history_id = $1", searchHistoryID)
	if err != nil {
		return entity.SearchResponse{}, err
	}

	for rows.Next() {
		var address entity.Address
		err = rows.Scan(&address.Lat, &address.Lng)
		if err != nil {
			return entity.SearchResponse{}, err
		}
		searchResponse.Addresses = append(searchResponse.Addresses, &address)
	}

	return searchResponse, nil
}
