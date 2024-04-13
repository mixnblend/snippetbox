package models

import (
	"testing"

	"github.com/mixnblend/snippetbox/internal/assert"
)

func TestUserModelExistsIntegration(t *testing.T) {
	integrationTest(t)

	testCases := []struct {
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

	for _, tableTest := range testCases {
		t.Run(tableTest.name, func(t *testing.T) {
			// given ... we have a database with a user record in it
			db := newTestDB(t)
			m := UserModel{DB: db}
			// when ... we try to check if a user exists
			result, err := m.Exists(tableTest.userID)
			// then ... the result should be true if it does and false if it doesn't
			assert.Equal(t, result, tableTest.want)
			assert.NilError(t, err)
		})
	}
}
