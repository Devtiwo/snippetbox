package models

import (
  "database/sql"
  "errors"
  "time"
)

// Defining a Snippet struct to hold the data for an individual snippet.
type Snippet struct {
  ID int
  Title string
  Content string
  Created time.Time
  Expires time.Time
}

// Defining a snippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
  DB *sql.DB
}

// Inserts a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
  // writing the SQL statement for inserting a new record into the snippets table.
  stmt := "INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))"
  // Using the EXEC() method on the embedded DB connection pool to execute the sql statement.
  result, err := m.DB.Exec(stmt, title, content, expires)
  if err != nil {
	return 0, err
  }
  // Using the LastInsertId() method on the result variable to get the ID of the newly inserted record.
  id, err := result.LastInsertId()
  if err != nil {
	return 0, err
  }
  return int(id), nil
}

// Retrieves a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
  stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?"
  row := m.DB.QueryRow(stmt, id)
  s := &Snippet{}
  err := row.Scan(&s.ID, &s.Title,&s.Content, &s.Created, &s.Expires)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
	  return nil, ErrNoRecord
	} else {
	  return nil, err
	}
  }
  return s, nil
}

// Returns the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
  stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10"
  rows, err := m.DB.Query(stmt)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  snippets := []*Snippet{}
  for rows.Next() {
	s := &Snippet{}
	err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
	  return nil, err
	}
	snippets = append(snippets, s)
  }
  if err = rows.Err(); err != nil {
	return nil, err
  }
  return snippets, nil
}