package mysql

import (
	"context"
	"database/sql"

	"github.com/k3forx/clean-architecture-with-Golang/domain"
)

type mysqlAuthorRepo struct {
	DB *sql.DB
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewMysqlAuthorRepository(db *sql.DB) domain.AuthorRepository {
	return &mysqlAuthorRepo{
		DB: db,
	}
}

func (m *mysqlAuthorRepo) getOne(ctx context.Context, query string, args ...interface{}) (res domain.Author, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Author{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.Author{}

	err = row.Scan(
		&res.Id,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}

func (m *mysqlAuthorRepo) GetByID(ctx context.Context, id int64) (domain.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(ctx, query, id)
}
