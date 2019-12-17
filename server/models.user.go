// models.user.go

package main

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	//"log"
	//"fmt"
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
	Login          time.Time
	CreatedOn      int64         `json:"created_on" bson:"created_on"`
	Level          int64         `json:"level" bson:"level"`
	Resume         bson.ObjectId `json:"resume,omitempty" bson:"resume,omitempty"`
}

type user_account struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username     string        `json:"username" form:"username" binding:"required" bson:"username"`
	Nickname     string        `json:"nickname" form:"nickname" binding:"required" bson:"nickname"`
	Password     string        `json:"password" form:"password" binding:"required" bson:"password"`
	Email        string        `json:"email" form:"email" binding:"required" bson:"email"`
	Status       string        `json:"status" form:"status" binding:"required" bson:"status"`
	Salt         string        `json:"salt" form:"salt" binding:"required" bson:"salt"`
	LastActivity time.Time
	CreatedOn    time.Time     `json:"created_on" bson:"created_on"`
	Level        int64         `json:"level" bson:"level"`
	UserGroupId  bson.ObjectId `json:"userGroupId" bson:"userGroupId,omitempty"`
	ProfileId    bson.ObjectId `json:"profileId" bson:"profileId,omitempty"`
}

type user_group struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name string
}

type user_profile struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	UserId   bson.ObjectId `json:"userId" bson:"userId"`
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

func UserSessionQuery(cws *careWorkerServer, username string) *user {
	queryUser := new(user)
	err := cws.collection["users"].Find(bson.M{"username": username}).One(&queryUser)
	if err != nil {
		panic(err)
	}
	dbgMessage("queryUser.Email:%s, Password:%s\n", queryUser.Email, queryUser.Password)

	return queryUser
}

func isUserValid(cws *careWorkerServer, email, password string) *user_account {
	queryUser := new(user_account)
	err := cws.collection["user_account"].Find(bson.M{"email": email}).One(&queryUser)
	if err != nil {
		panic(err)
	}
	dbgMessage("queryUser.Email:%s, Password:%s\n", queryUser.Email, queryUser.Password)

	if queryUser.Email == email && queryUser.Password == password {
		return queryUser
	}
	return nil
}

func registerNewUser(cws *careWorkerServer, newUserAccount *user_account) (*user_account, error) {
	newUserProfile := new(user_profile)

	if strings.TrimSpace(newUserAccount.Password) == "" {
		dbgMessage("registerNewUserAccount password null\n")
		return nil, errors.New("The password can't be empty")
	} else if !isUserEmailAvailable(cws, newUserAccount.Email) {
		dbgMessage("registerNewUserAccount email exlist\n")
		return nil, errors.New("The email isn't available")
	}

	newUserAccount.Id = bson.NewObjectId()
	err := cws.collection["user_account"].Insert(newUserAccount)
	if err != nil {
		panic(err)
	}

	newUserProfile.UserId = newUserAccount.Id
	err = cws.collection["user_profile"].Insert(newUserProfile)
	if err != nil {
		panic(err)
	}
	dbgMessage("registerNewUser success\n")
	return newUserAccount, nil
}

func isUserEmailAvailable(cws *careWorkerServer, email string) bool {
	result := user_account{}
	cws.collection["user_account"].Find(bson.M{"email": email}).One(&result)
	if result.Email == email {
		return false
	}
	return true
}

func isUserSaltAvailable(cws *careWorkerServer, email string) (*user_account, bool) {
	saltUser := new(user_account)
	err := cws.collection["user_account"].Find(bson.M{"email": email}).One(&saltUser)
	if err != nil {
		panic(err)
	}
	if saltUser.Email == email {
		return saltUser, true
	}
	return nil, false
}

func updateUserProfile(cws *careWorkerServer, userId string, profile *user_profile) error {
	userProfile := new(user_profile)

	err := cws.collection["user_profile"].Find(bson.M{"userId": userId}).One(&userProfile)
	if err != nil {
		dbgMessage("insert profile:%s", userProfile.Id.Hex())
		err = cws.collection["user_profile"].Insert(profile)
	} else {
		dbgMessage("update profile:%s", userProfile.Id.Hex())
		err = cws.collection["user_profile"].Update(bson.M{"_id": userProfile.Id}, profile)
	}
	return nil
}
