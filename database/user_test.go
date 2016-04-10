package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterUserCreatesUserInCorrectState(t *testing.T) {

}

func TestNewSessionCreatesSession(t *testing.T) {
	assert := assert.New(t)

	user := TestUser([]string{})
	cookie, err := user.NewSession()
	assert.Nil(err)
	session := Session{}
	DB().Where(Session{Cookie: cookie}).First(&session)
	if assert.NotNil(session.UserID) {
		assert.Equal(session.UserID, user.ID)
	}
}