package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Client struct {
	db *sql.DB
}

var (
	sharedClient *Client
	mu           sync.Mutex
)

var ErrClientNotInitialized = errors.New("database not initialized")

func sanitizeDatabaseURL(databaseURL string) (string, error) {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return "", err
	}

	q := parsed.Query()
	if schema := q.Get("schema"); schema != "" {
		q.Del("schema")
		if q.Get("search_path") == "" {
			q.Set("search_path", schema)
		}
		parsed.RawQuery = q.Encode()
	}

	return parsed.String(), nil
}

func NewClient(databaseURL string) (*Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if sharedClient != nil {
		return sharedClient, nil
	}

	sanitizedURL, err := sanitizeDatabaseURL(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse postgres url: %w", err)
	}

	sqlDB, err := sql.Open("pgx", sanitizedURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	sharedClient = &Client{db: sqlDB}
	return sharedClient, nil
}

func GetClient() (*Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if sharedClient == nil {
		return nil, ErrClientNotInitialized
	}

	return sharedClient, nil
}

func (c *Client) DB() *sql.DB {
	return c.db
}

func (c *Client) Disconnect() {
	mu.Lock()
	defer mu.Unlock()

	if sharedClient == nil {
		return
	}

	_ = sharedClient.db.Close()
	sharedClient = nil
}
