package postgres

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func isAlreadyExistsError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func queryBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func queryLike(q squirrel.SelectBuilder, column, term string) squirrel.SelectBuilder {
	return q.Where(column+" LIKE ?", "%"+term+"%")
}
