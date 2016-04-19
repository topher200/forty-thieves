package dal

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func (u *User) signupNewUserRowForTest(t *testing.T) *UserRow {
	userRow, err := u.Signup(nil, newEmailForTest(), "abc123", "abc123")
	assert.NotNil(t, userRow)
	assert.Nil(t, err)
	assert.True(t, userRow.ID > 0, "user should be given a real ID")
	return userRow
}

func TestUserCRUD(t *testing.T) {
	u := NewUserDBForTest(t)

	// Signup
	userRow := u.signupNewUserRowForTest(t)

	// DELETE FROM users WHERE id=...
	_, err := u.DeleteById(nil, userRow.ID)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}
