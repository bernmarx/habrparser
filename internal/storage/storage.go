//go:generate mockgen -source $GOFILE -destination ./storage_mock.go -package $GOPACKAGE
package storage

import (
	"database/sql"
	"log"
	"os"

	"github.com/bernmarx/habrparser/internal/page"
)

type storage interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Storage struct {
	storage
}

func NewStorage() (*Storage, error) {
	connData := "host=" + os.Getenv("DB_HOST") + " " + "port=" + os.Getenv("DB_PORT")
	connData = connData + " " + "user=" + os.Getenv("DB_USER") + " " + "password=" + os.Getenv("DB_PASSWORD")
	connData = connData + " " + "dbname=" + os.Getenv("DB_NAME") + " " + "sslmode=" + os.Getenv("DB_SSLMODE")
	log.Println(connData)

	db, err := sql.Open("postgres", connData)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	s := Storage{storage: db}

	return &s, nil
}

func (s *Storage) AddPageData(p *page.Page) error {
	sqlstmt1 := `SELECT addpage($1)`
	sqlstmt2 := `SELECT addpagejson($1, $2)`
	pageJSON := p.GetJSON()

	_, err := s.Exec(sqlstmt1, pageJSON)
	if err != nil {
		return err
	}

	_, err = s.Exec(sqlstmt2, p.ID, pageJSON)
	if err != nil {
		return err
	}

	return nil
}
