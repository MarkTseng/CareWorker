// models.user.go

package main

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"strings"
	"time"
)

type user_account struct {
	Id           bson.ObjectId `json:"_id,omitempty" bson:"id,omitempty"`
	Password     string
	Email        string
	Created      time.Time
	LastActivity time.Time
	Status       string
	Level        int
	Salt         string        `json:"salt" form:"salt" binding:"required" bson:"salt"`
	UserGroupId  bson.ObjectId `json:"_id,omitempty" bson:"id,omitempty"`
	ResumeId     bson.ObjectId `json:"_id,omitempty" bson:"id,omitempty"`
}

type user_group struct {
	Id   bson.ObjectId `json:"_id,omitempty" bson:"id,omitempty"`
	Name string
}

type user_profile struct {
	Id       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name" form:"name" binding:"required" bson:"name"`
	Birthday string        `json:"birthday" form:"birthday" binding:"required" bson:"birthday"`
	Phone    string        `json:"phone" form:"phone" binding:"required" bson:"phone"`
	City     string        `json:"city" form:"city" binding:"required" bson:"city"`
	District string        `json:"district" form:"district" binding:"required" bson:"district"`
	Gender   string        `json:"gender" form:"gender" binding:"required" bson:"gender"`
	Street   string        `json:"street" form:"street" binding:"required" bson:"street"`
	Zipcode  string        `json:"zipcode" form:"zipcode" binding:"required" bson:"zipcode"`
	JobsBrew string        `json:"jobbrew" form:"jobbrew" binding:"required" bson:"jobbrew"`
	License  string        `json:"license" form:"license" binding:"required" bson:"license"`
}

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
	Login          time.Time
	CreatedOn      int64         `json:"created_on" bson:"created_on"`
	Level          int64         `json:"level" bson:"level"`
	Resume         bson.ObjectId `json:"resume,omitempty" bson:"resume,omitempty"`
}

func UserSessionQuery(cws *careWorkerServer, username string) *user {
	queryUser := new(user)
	cws.collection["users"].Find(bson.M{"username": username}).One(&queryUser)
	dbgMessage("queryUser.Email:%s, Password:%s\n", queryUser.Email, queryUser.Password)

	return queryUser
}

func isUserValid(cws *careWorkerServer, email, password string) *user {
	queryUser := new(user)
	cws.collection["users"].Find(bson.M{"email": email}).One(&queryUser)
	dbgMessage("queryUser.Email:%s, Password:%s\n", queryUser.Email, queryUser.Password)

	if queryUser.Email == email && queryUser.Password == password {
		return queryUser
	}
	return nil
}

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

func isUserEmailAvailable(cws *careWorkerServer, email string) bool {
	result := user{}
	cws.collection["users"].Find(bson.M{"email": email}).One(&result)
	if result.Email == email {
		return false
	}
	return true
}

func isUserSaltAvailable(cws *careWorkerServer, email string) (*user, bool) {
	saltUser := new(user)
	cws.collection["users"].Find(bson.M{"email": email}).One(&saltUser)
	if saltUser.Email == email {
		return saltUser, true
	}
	return nil, false
}

func updateUserProfile(cws *careWorkerServer, email string, profile *user_profile) error {
	var err error
	userProfile := new(user_profile)

	err = cws.collection["user_profile"].Find(bson.M{"name": email}).One(&userProfile)
	dbgMessage("userProfile:%s", userProfile.Name)
	if err != nil {
		dbgMessage("insert:%s", userProfile.Name)
		err = cws.collection["user_profile"].Insert(profile)
	} else {
		dbgMessage("update:%s", userProfile.Name)
		err = cws.collection["user_profile"].Update(bson.M{"_id": userProfile.Id}, profile)
	}
	return nil
}
