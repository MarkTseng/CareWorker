// models.user.go

package main

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

/*
// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}
*/
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

	if result.Username == username && result.Password == password {
		return true
	}
	return false
}

// Register a new user with the given username and password
// NOTE: For this demo, we
func registerNewUser(cws *careWorkerServer, username, password, salt string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		log.Printf("registerNewUser password null\n")
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(cws, username) {
		log.Printf("registerNewUser username exlist\n")
		return nil, errors.New("The username isn't available")
	}

	u := user{Username: username, Password: password, Salt: salt}
	cws.users.Insert(&u)

	log.Printf("registerNewUser success\n")
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

// Check if the supplied username is available
func isUserSaleAvailable(cws *careWorkerServer, username string) (*user, bool) {
	result := new(user)
	cws.users.Find(bson.M{"username": username}).One(&result)
	if result.Username == username {
		return result, true
	}
	return nil, false
}
