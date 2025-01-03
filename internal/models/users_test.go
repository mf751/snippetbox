package models

import (
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v4/stdlib"

	"github.com/mf751/snippetbox/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	// Register it only once before all the tests to avoid the registering twice error
	sql.Register("postgres", stdlib.GetDefaultDriver())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)

			model := UserModel{db}

			exists, err := model.Exists(test.userID)
			assert.Equal(t, exists, test.want)
			assert.NilError(t, err)
		})
	}
}
