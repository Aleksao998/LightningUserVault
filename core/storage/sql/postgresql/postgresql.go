package postgresql

import (
	"database/sql"
	"errors"

	"github.com/Aleksao998/LightingUserVault/core/common"
	_ "github.com/lib/pq" //nolint: blank
)

var errUserNotFound = errors.New("user not found")

const (
	getUserQuery = `SELECT id, name FROM users WHERE id = $1`
	addUserQuery = `INSERT INTO users (name) VALUES ($1) RETURNING id`
)

type Storage struct {
	db          *sql.DB
	getUserStmt *sql.Stmt
	addUserStmt *sql.Stmt
}

// NewStorage initializes a new Storage instance with a database at the given path.
func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	getUserStmt, err := db.Prepare(getUserQuery)
	if err != nil {
		return nil, db.Close()
	}

	addUserStmt, err := db.Prepare(addUserQuery)
	if err != nil {
		return nil, db.Close()
	}

	return &Storage{
		db:          db,
		getUserStmt: getUserStmt,
		addUserStmt: addUserStmt,
	}, nil
}

// Get retrieves the value for a given key and returns an error if any issue occurs during the operation.
func (p *Storage) Get(id int) (*common.User, error) {
	user := &common.User{}

	err := p.getUserStmt.QueryRow(id).Scan(&user.ID, &user.Name)

	if err == sql.ErrNoRows {
		return nil, errUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Set stores a value for a given key and returns an error if any issue occurs during the operation.
func (p *Storage) Set(name string) (int64, error) {
	var id int64

	err := p.addUserStmt.QueryRow(name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Close closes the database connection and returns an error if any issue occurs during the operation.
func (p *Storage) Close() error {
	return p.db.Close()
}
