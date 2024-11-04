package models

import (
	"database/sql"
	"errors"
	"time"
)

// db entity
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// model/repo/data access layer/dao
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `insert into snippets (title, content, created, expires)
	values(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	r, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	s := &Snippet{}
	stmt := `select * from snippets where expires > UTC_TIMESTAMP() and id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// returning our own sentinel error to abstract the datastore specific errors.
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// returns 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {

	snippets := []*Snippet{}
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// to close the connections from db if methods fails
	defer rows.Close()

	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// To check if rows.Next() ends because of an error
	// or if no next row is found
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

// upside of writing all the code of sql - like connecting to db
// is the it's non-magical and we can understand and
// control exactly what is going on

// extension of db/sql pkg : https://jmoiron.github.io/sqlx/
// transactions
