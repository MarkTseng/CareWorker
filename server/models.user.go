// models.user.go

package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	//"log"
	//"fmt"
	"strings"
	"time"
)

type user_account struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username    string        `json:"username" form:"username" binding:"required" bson:"username"`
	Nickname    string        `json:"nickname" form:"nickname" binding:"required" bson:"nickname"`
	Password    string        `json:"password" form:"password" binding:"required" bson:"password"`
	Email       string        `json:"email" form:"email" binding:"required" bson:"email"`
	Status      string        `json:"status" form:"status" binding:"required" bson:"status"`
	Salt        string        `json:"salt" form:"salt" binding:"required" bson:"salt"`
	ResetCode   string        `json:"resetcode" form:"resetcode" binding:"required" bson:"resetcode"`
	CreatedOn   time.Time     `json:"created_on" bson:"created_on"`
	VIPEndTime  time.Time     `json:"vip_endtime" bson:"vip_endtime"`
	Level       int64         `json:"level" bson:"level"`
	UserGroupId bson.ObjectId `json:"userGroupId" bson:"userGroupId,omitempty"`
	ProfileId   bson.ObjectId `json:"profileId" bson:"profileId,omitempty"`
}

type user_group struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name string
}

type user_profile struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	UserId   bson.ObjectId `json:"userId" bson:"userId"`
	IdType   string        `json:"idtype" form:"idtype" binding:"required" bson:"idtype"`
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

func isUserValid(cws *careWorkerServer, email, password string) *user_account {
	queryUser := new(user_account)
	err := cws.collection["user_account"].Find(bson.M{"email": email}).One(&queryUser)
	if err != nil {
		//panic(err)
		return nil
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
		//panic(err)
		return nil, false
	}

	if saltUser.Email == email {
		return saltUser, true
	}
	return nil, false
}

func updateUserProfile(cws *careWorkerServer, userId bson.ObjectId, profile *user_profile) error {
	userProfile := new(user_profile)

	err := cws.collection["user_profile"].Find(bson.M{"userId": userId}).One(&userProfile)
	if err != nil {
		panic(err)
	}
	dbgMessage("update profile:%s", userProfile.Id.Hex())
	err = cws.collection["user_profile"].Update(bson.M{"_id": userProfile.Id}, profile)

	return nil
}

func getUserProfile(cws *careWorkerServer, userId string) ([]user_profile, error) {
	userProfile := user_profile{}

	dbgMessage("getUserProfile: %s", userId)

	err := cws.collection["user_profile"].Find(bson.M{"userId": bson.ObjectIdHex(userId)}).One(&userProfile)
	if err != nil {
		panic(err)
	}

	userProfileResult := []user_profile{userProfile}
	return userProfileResult, nil
}

func GenerateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	dbgMessage("Hash to store: %s\n", string(hash))

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func setResetCode(cws *careWorkerServer, email string) (string, bool) {
	userAccount := new(user_account)

	cws.collection["user_account"].Find(bson.M{"email": email}).One(&userAccount)

	// get resetcode
	userAccount.ResetCode = GenerateToken(email)

	// set resetcode
	err := cws.collection["user_account"].Update(bson.M{"email": userAccount.Email}, userAccount)
	if err != nil {
		//panic(err)
		return "", false
	}

	return userAccount.ResetCode, true
}

func clearResetCode(cws *careWorkerServer, email string) (string, bool) {
	userAccount := new(user_account)

	cws.collection["user_account"].Find(bson.M{"email": email}).One(&userAccount)

	// get resetcode
	userAccount.ResetCode = ""

	// set resetcode
	err := cws.collection["user_account"].Update(bson.M{"email": userAccount.Email}, userAccount)
	if err != nil {
		panic(err)
		return "", false
	}

	return userAccount.ResetCode, true
}

func verifyResetCode(cws *careWorkerServer, email string, resetCode string) bool {
	userAccount := new(user_account)

	err := cws.collection["user_account"].Find(bson.M{"email": email}).One(&userAccount)

	if err != nil {
		panic(err)
		return false
	}

	if userAccount.ResetCode == resetCode && strings.TrimSpace(userAccount.ResetCode) != "" {
		return true
	}

	return false
}

func resetPassword(cws *careWorkerServer, email string, password string) (string, bool) {
	userAccount := new(user_account)

	cws.collection["user_account"].Find(bson.M{"email": email}).One(&userAccount)

	userAccount.Password = password

	// set resetcode
	err := cws.collection["user_account"].Update(bson.M{"email": userAccount.Email}, userAccount)
	if err != nil {
		panic(err)
		return "", false
	}

	return userAccount.ResetCode, true
}
