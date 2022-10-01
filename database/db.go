package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	queryCreateUsers = `
    CREATE TABLE users (
        id INTEGER PRIMARY KEY,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME
    );`
)

var store *Store

func GetStore() *Store {
	if store != nil {
		return store
	}

	var err error
	store, err = newStore()
	if err != nil {
		log.Fatal(err)
	}
	return store
}

type Store struct {
	db *sql.DB
}

func newStore() (*Store, error) {

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		return nil, err
	}

	if db.Ping() != nil {
		return nil, err
	}

	_, err = db.Exec(queryCreateUsers)
	if err != nil && err.Error() != "table users already exists" {
		return nil, err
	}

	return &Store{db}, nil
}

func (s *Store) CreateUser(username, password string) (int64, error) {
	createdAt := time.Now()

	res, err := s.db.Exec(
		`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`,
		username, password, createdAt,
	)

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (s *Store) Users() ([]User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	var users []User

	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	//fmt.Printf("%#v", users)
	rows.Close()

	return users, nil
}

func (s *Store) UserById(userId int64) (*User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE id=$1", userId)

	var user User
	switch err := row.Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("user with id: %d not found", userId)
	case nil:
		return &user, nil
	default:
		return nil, err
	}
}

func (s *Store) Close() error {
	return s.db.Close()
}
