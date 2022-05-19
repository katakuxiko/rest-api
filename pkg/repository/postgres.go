package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
	todoListTable = "todo_lists"
	userListTable = "users_lists"
	todoItemsTable = "todo_items"
	listItemsTable = "lists_items"
)
type Config struct {
	Host string 
	Port string
	Username string
	Password string
	DBname string
	SSLmode string
}

func NewPostgresDb(cfg Config)(*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",cfg.Host,cfg.Port,cfg.Username,cfg.DBname, cfg.Password,cfg.SSLmode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil

}