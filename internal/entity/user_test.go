package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user, err := NewUser("Felipe Dias", "felipe@gmail.com", "123456")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "Felipe Dias", user.Name)
	assert.Equal(t, "felipe@gmail.com", user.Email)
}

func TestUser_ComparePassword(t *testing.T) {
	user, err := NewUser("Felipe Dias", "felipe@gmail.com", "123456")
	assert.Nil(t, err)
	assert.True(t, user.ComparePassword("123456"))
	assert.False(t, user.ComparePassword("1234567"))
	assert.NotEqual(t, "123456", user.Password)
}
