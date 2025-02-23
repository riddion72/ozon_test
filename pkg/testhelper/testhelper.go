package testhelper

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

type SetupDBMockOption func(testing.TB, sqlxmock.Sqlmock)

func SetupDBMock(tb testing.TB, opts ...SetupDBMockOption) *sqlx.DB {
	tb.Helper()

	db, dbMock, err := sqlxmock.Newx(sqlxmock.QueryMatcherOption(sqlxmock.QueryMatcherEqual))
	require.NoError(tb, err)

	tb.Cleanup(func() {
		assert.NoError(tb, dbMock.ExpectationsWereMet())
	})

	for _, o := range opts {
		o(tb, dbMock)
	}

	return db
}
