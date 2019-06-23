package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DBConnector struct {
	db *sql.DB
	c  *Client
}

func NewConnector(db *sql.DB, c *Client) (*DBConnector, error) {
	dc := &DBConnector{
		db: db,
		c:  c,
	}

	return dc, nil
}

//
func (dc *DBConnector) AddSearchQueue(ids []string) error {
	return nil

}

func (dc *DBConnector) Close() error {
	return nil
}
