package mysql

import (
	"database/sql"

	"github.com/FollowTheProcess/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// Get a snippet by id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Latest will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
