package imglib

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestUserCreation(t *testing.T) {
	db := openTempDB()

	//given: a user object
	user := user()

	//when: I create the user
	err := db.CreateUser(&user)
	assert.NoError(t, err)

	//then: The autoset fields are set
	assert.NotNil(t, user.ID)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
	assert.Nil(t, user.DeletedAt)

	//then: I can retrieve it again
	userNew, err := db.UserByUsername("fmuell")
	assert.NoError(t, err)
	assertUsersAreEqual(t, &user, userNew)
}

func TestUserUpdate(t *testing.T) {
	db := openTempDB()

	//given: a user object in db
	user := user()
	err := db.CreateUser(&user)
	assert.NoError(t, err)
	assert.NotNil(t, user.ID)

	//when I change the user
	user.NickName = "Foo"
	errUpdate := db.SaveUser(&user)
	assert.NoError(t, errUpdate)

	//then: I can retrieve it again
	userNew, err := db.UserByUsername("fmuell")
	assert.NoError(t, err)
	assertUsersAreEqual(t, &user, userNew)
}

func user() User {
	return User{
		UserName: "fmuell",
		NickName: "Franky",
		Link:     "http://www.example.com/frank/muelller"}
}

func assertUsersAreEqual(t *testing.T, user *User, userNew *User) {
	assert.Equal(t, user.ID, userNew.ID)
	assert.Equal(t, user.UserName, userNew.UserName)
	assert.Equal(t, user.NickName, userNew.NickName)
	assert.Equal(t, user.Link, userNew.Link)
	assert.True(t, user.CreatedAt.Equal(userNew.CreatedAt))
	assert.True(t, user.UpdatedAt.Equal(userNew.UpdatedAt))
}

func openTempDB() *ImageLibrary {
	file, _ := ioutil.TempDir("", "galleryDir.")

	db := &ImageLibrary{}
	if err := db.Open(file); err != nil {
		log.Fatal(err.Error())
	}
	return db
}
