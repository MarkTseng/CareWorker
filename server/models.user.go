// models.user.go

package main

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"strings"
	"time"
)

type user struct {
	Id             bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Username       string        `json:"username" form:"username" binding:"required" bson:"username"`
	Nickname       string        `json:"nickname" form:"nickname" binding:"required" bson:"nickname"`
	Password       string        `json:"password" form:"password" binding:"required" bson:"password"`
	Email          string        `json:"email" form:"email" binding:"required" bson:"email"`
	Birthday       string        `json:"birthday" form:"birthday" binding:"required" bson:"birthday"`
	CellPhone      string        `json:"cellPhone" form:"cellPhone" binding:"required" bson:"cellPhone"`
	City           string        `json:"city" form:"city" binding:"required" bson:"city"`
	District       string        `json:"district" form:"district" binding:"required" bson:"district"`
	Gender         string        `json:"gender" form:"gender" binding:"required" bson:"gender"`
	Idtype         string        `json:"idtype" form:"idtype" binding:"required" bson:"idtype"`
	RequestService []string      `json:"requestService" form:"requestService" binding:"required" bson:"requestService"`
	Contact        []string      `json:"contact" form:"contact" binding:"required" bson:"contact"`
	Street         string        `json:"street" form:"street" binding:"required" bson:"street"`
	Zipcode        string        `json:"zipcode" form:"zipcode" binding:"required" bson:"zipcode"`
	Salt           string        `json:"salt" form:"salt" binding:"required" bson:"salt"`
	login          time.Time
}

// Check if the username and password combination is valid
func isUserValid(cws *careWorkerServer, email, password string) *user {
	queryUser := new(user)
	cws.collection["users"].Find(bson.M{"email": email}).One(&queryUser)
	dbgMessage("queryUser.Email:%s, Password:%s\n", queryUser.Email, queryUser.Password)

	if queryUser.Email == email && queryUser.Password == password {
		return queryUser
	}
	return nil
}

// Register a new user with the given username and password
// NOTE: For this demo, we
func registerNewUser(cws *careWorkerServer, newUser *user) (*user, error) {
	if strings.TrimSpace(newUser.Password) == "" {
		dbgMessage("registerNewUser password null\n")
		return nil, errors.New("The password can't be empty")
	} else if !isUserEmailAvailable(cws, newUser.Email) {
		dbgMessage("registerNewUser email exlist\n")
		return nil, errors.New("The email isn't available")
	}

	cws.collection["users"].Insert(newUser)
	dbgMessage("registerNewUser success\n")
	return newUser, nil
}

// Check if the supplied username is available
func isUserEmailAvailable(cws *careWorkerServer, email string) bool {
	result := user{}
	cws.collection["users"].Find(bson.M{"email": email}).One(&result)
	if result.Email == email {
		return false
	}
	return true
}

// Check if the supplied username is available
func isUserSaleAvailable(cws *careWorkerServer, email string) (*user, bool) {
	saltUser := new(user)
	cws.collection["users"].Find(bson.M{"email": email}).One(&saltUser)
	if saltUser.Email == email {
		return saltUser, true
	}
	return nil, false
}
