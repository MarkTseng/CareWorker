// models.user.go

package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

const (
	MongoDBUrl = "mongodb://localhost:27017/articles_demo_dev"
)

type user struct {
	Id       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string        `json:"username" form:"username" binding:"required" bson:"username"`
	Password string        `json:"password" form:"password" binding:"required" bson:"password"`
}

// Check if the username and password combination is valid
func isUserValid(username, password string) bool {
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("users")
	result := user{}
	err = c.Find(bson.M{"username": username}).One(&result)

	if result.Username == username && result.Password == password {
		return true
	}
	return false
}

// Register a new user with the given username and password
// NOTE: For this demo, we
func registerNewUser(username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	u := user{Username: username, Password: password}

	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("users")
	err = c.Insert(&u)

	return &u, nil
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("users")
	result := user{}
	err = c.Find(bson.M{"username": username}).One(&result)

	if result.Username == username {
		return false
	}
	return true
}
