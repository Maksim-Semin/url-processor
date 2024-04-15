package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	lg "main/pkg/logger"
	"main/resource"
)

type urlInfo struct {
	shortCode string
	URL       string
}

var conn = "host=db port=5432 user=postgres dbname=urls password=password sslmode=disable"

type DB struct {
	db *sql.DB
}

var Db *DB

func init() {
	if resource.CFG.DB == "postgres" {
		Db = NewPostgresDB()
	}
}

func NewPostgresDB() *DB {
	return &DB{}
}

func (d *DB) OpenConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		lg.Logger.Error("error on connection to DB", zap.Error(err))
		return db, err
	}
	return db, nil
}

func (d *DB) SaveUrl(shortcode, URL string) (string, error) {
	db, err := d.OpenConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var existingShortcode string
	err = db.QueryRow(`SELECT shortCode FROM URLs WHERE URL = $1`, URL).Scan(&existingShortcode)
	if err == nil {
		lg.Logger.Warn("duplicate short code for url", zap.String("url", URL))
		return existingShortcode, nil
	} else if err != sql.ErrNoRows {
		lg.Logger.Error("error on executing query to DB", zap.Error(err))
		return "", err
	}

	_, err = db.Exec(`INSERT INTO URLs (shortCode, URL) VALUES ($1, $2)`, shortcode, URL)
	if err != nil {
		lg.Logger.Error("error on executing query to DB", zap.Error(err))
		return "", err
	}

	lg.Logger.Info("data was successfully saved", zap.String("shortcode", shortcode), zap.String("url", URL))
	return shortcode, nil
}

func (d *DB) GetURL(shortcode string) (string, error) {
	var info urlInfo
	db, err := d.OpenConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	row, err := db.Query(`SELECT URL FROM URLs WHERE shortCode = $1`, shortcode)
	if err != nil {
		lg.Logger.Error("error on executing query to DB", zap.Error(err))
		return "", err
	}
	defer row.Close()

	if !row.Next() {
		return "", sql.ErrNoRows
	}

	err = row.Scan(&info.URL)
	if err != nil {
		lg.Logger.Error("error on scan query from DB", zap.Error(err))
		return "", err
	}
	return info.URL, nil
}
