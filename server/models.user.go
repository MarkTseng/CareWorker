// models.user.go

package main

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type user struct {
	Id       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string        `json:"username" form:"username" binding:"required" bson:"username"`
	Password string        `json:"password" form:"password" binding:"required" bson:"password"`
	Salt     string        `json:"salt" form:"salt" binding:"required" bson:"salt"`
	login    time.Time
}

// Check if the username and password combination is valid
func isUserValid(cws *careWorkerServer, username, password string) bool {
	result := new(user)
	cws.users.Find(bson.M{"username": username}).One(&result)

	// use password from http parameter
	salt_password := DoHash(password, result.Salt)

	if result.Username == username && result.Password == salt_password {
		return true
	}
	return false
}

// Register a new user with the given username and password
// NOTE: For this demo, we
func registerNewUser(cws *careWorkerServer, username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(cws, username) {
		return nil, errors.New("The username isn't available")
	}

	// generate salt
	salt := genSaltString()
	// get salt password with sha256 hexcode
	salt_pass := DoHash(password, salt)
	u := user{Username: username, Password: salt_pass, Salt: salt}

	cws.users.Insert(&u)

	return &u, nil
}

// Check if the supplied username is available
func isUsernameAvailable(cws *careWorkerServer, username string) bool {
	result := user{}
	cws.users.Find(bson.M{"username": username}).One(&result)
	if result.Username == username {
		return false
	}
	return true
}
