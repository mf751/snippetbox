package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Becasue the application struct is to expect this
// so that it can recive both the snippetModel and the mocks.SnippetModel
type SnippetModelInterface interface {
	Insert(title, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Insert(title, content string, expires int) (int, error) {
	sqlStatment := `INSERT INTO snippets (title, content, created, expires)
VALUES($1, $2, NOW(), NOW() + $3 * INTERVAL '1 DAY') RETURNING id;`

	var lastId int
	err := model.DB.QueryRow(sqlStatment, title, content, expires).Scan(&lastId)
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}

func (model *SnippetModel) Get(id int) (*Snippet, error) {
	sqlStatment := `SELECT * FROM snippets WHERE expires > NOW() AND id = $1`
	snippet := &Snippet{}
	err := model.DB.QueryRow(sqlStatment, id).Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return snippet, nil
}

func (model *SnippetModel) Latest() ([]*Snippet, error) {
	sqlStatment := `SELECT * FROM snippets WHERE expires > NOW() ORDER BY id DESC LIMIT 10`
	rows, err := model.DB.Query(sqlStatment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}
	for rows.Next() {
		snippet := &Snippet{}
		err = rows.Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Created,
			&snippet.Expires,
		)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
