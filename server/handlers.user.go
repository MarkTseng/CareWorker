// handlers.user.go
package main

import (
	//"math/rand"
	"net/http"
	//"strconv"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	//"log"
	"fmt"
)

type responMSG struct {
	Message string `json:"Message"`
}

func generateSessionToken() string {
	return uuid.New()
}

func (cws *careWorkerServer) getProfile(c *gin.Context) {
	userId := c.Param("userId")

	if userProfile, err := getUserProfile(cws, userId); err == nil {
		c.SecureJSON(http.StatusOK, userProfile)

	} else {
		ErrMSG := []responMSG{{Message: "get user profile fail"}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
	}
}

func (cws *careWorkerServer) profile(c *gin.Context) {
	profile := new(user_profile)
	// Obtain the POSTed JSON username and password values
	if err := c.ShouldBindJSON(&profile); err != nil {
		ErrMSG := []responMSG{{Message: err.Error()}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
		return
	}
	dbgMessage("profile JSON name:%s, phone:%s\n", profile.Id.Hex(), profile.Phone)
	fmt.Println(profile)

	if err := updateUserProfile(cws, profile.UserId, profile); err == nil {
		dbgMessage("%s: register success\n", profile.Id.Hex())
		successMSG := []responMSG{{Message: "Success"}}
		c.SecureJSON(http.StatusOK, successMSG)

	} else {
		ErrMSG := []responMSG{{Message: "update profile fail"}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
	}
}

type responUserInfo struct {
	Username string `json:"Username"`
	UserId   string `json:"UserId"`
}

func (cws *careWorkerServer) performLogin(c *gin.Context) {
	loginUserAccount := new(user_account)

	// Obtain the POSTed JSON username and password values
	if err := c.ShouldBindJSON(&loginUserAccount); err != nil {
		ErrMSG := []responMSG{{Message: err.Error()}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
		return
	}
	dbgMessage("performLogin JSON email:%s, password:%s\n", loginUserAccount.Email, loginUserAccount.Password)

	// Check if the username/password combination is valid
	if user := isUserValid(cws, loginUserAccount.Email, loginUserAccount.Password); user != nil {
		// If the username/password is valid set the token in a cookie
		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

		// save username in session
		session := sessions.Default(c)
		session.Set("username", user.Username)
		dbgMessage("set %s to session, token=%s", user.Username, token)
		err := session.Save()
		if err != nil {
			dbgMessage("user session svae failed")
		}
		UserInfo := []responUserInfo{{Username: user.Username, UserId: user.Id.Hex()}}
		dbgMessage("UserInfo:%s\n", UserInfo)
		c.SecureJSON(http.StatusOK, UserInfo)
	} else {
		ErrMSG := []responMSG{{Message: "login fail, Please chech account and password"}}
		c.SecureJSON(http.StatusUnauthorized, ErrMSG)
	}
}

func (cws *careWorkerServer) logout(c *gin.Context) {
	dbgMessage("logout\n")
	// Clear the cookie
	c.SetCookie("token", "", -1, "", "", false, true)
	session := sessions.Default(c)
	username := session.Get("username")
	if username != nil {
		//log.Printf("delete %s\n", username)
		session.Delete("username")
		session.Save()
	}
	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (cws *careWorkerServer) register(c *gin.Context) {
	newUserAccount := new(user_account)

	// Obtain the POSTed JSON username and password values
	if err := c.ShouldBindJSON(&newUserAccount); err != nil {
		dbgMessage("%s: register fail\n", newUserAccount.Email)
		ErrMSG := []responMSG{{Message: err.Error()}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
		return
	}

	if _, err := registerNewUser(cws, newUserAccount); err == nil {
		dbgMessage("%s: register success\n", newUserAccount.Email)
		successMSG := []responMSG{{Message: "Success"}}
		c.SecureJSON(http.StatusOK, successMSG)

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		ErrMSG := []responMSG{{Message: err.Error()}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
	}
}

type responUserSalt struct {
	Salt     string `json:"salt"`
	Username string `json:"username"`
}

func (cws *careWorkerServer) registerSalt(c *gin.Context) {
	// Obtain the POSTed JSON username and password values
	saltUser := new(user_account)
	if err := c.ShouldBindJSON(&saltUser); err != nil {
		ErrMSG := []responMSG{{Message: err.Error()}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
		return
	}
	dbgMessage("register JSON email:%s\n", saltUser.Email)

	if queryUser, err := isUserSaltAvailable(cws, saltUser.Email); err == true {
		UserSalt := []responUserSalt{{Salt: queryUser.Salt, Username: queryUser.Username}}
		c.SecureJSON(http.StatusOK, UserSalt)
	}
}

func (cws *careWorkerServer) forgotPassword(c *gin.Context) {
	userEmail := c.Param("email")

	dbgMessage("forgot email :%s\n", userEmail)
	// save reset hashcode in user database
	resetCode, err := setResetCode(cws, userEmail)

	if err == false {
		ErrMSG := []responMSG{{Message: "Wrong email address"}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
	} else {

		// send mail
		sendPasswordResetMail(userEmail, resetCode)

		successMSG := []responMSG{{Message: "Success"}}
		c.SecureJSON(http.StatusOK, successMSG)
	}
}

func (cws *careWorkerServer) resetPassword(c *gin.Context) {
	userEmail := c.Param("email")
	userResetCode := c.Param("resetcode")
	userNewPassword := c.Param("newpassword")

	// resetcode is correct
	dbgMessage("userEmail: %s, resetcode:%s, password:%s\n", userEmail, userResetCode, userNewPassword)

	// check resetcode flow
	ret := verifyResetCode(cws, userEmail, userResetCode)

	// resetcode is fail
	if ret == false {
		ErrMSG := []responMSG{{Message: "Verify Failed"}}
		c.SecureJSON(http.StatusBadRequest, ErrMSG)
	} else {
		// clear resetCode
		clearResetCode(cws, userEmail)
		// update password
		resetPassword(cws, userEmail, userNewPassword)
		successMSG := []responMSG{{Message: "Success"}}
		c.SecureJSON(http.StatusOK, successMSG)
	}
}
