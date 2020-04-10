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

func generateSessionToken() string {
	return uuid.New()
}

func (cws *careWorkerServer) getProfile(c *gin.Context) {
	profile := new(user_profile)
	userId := c.Param("userId")
	// Obtain the POSTed JSON username and password values
	//if err := c.ShouldBindJSON(&profile); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"payload": err.Error()})
	//	return
	//}
	dbgMessage("profile JSON name:%s, phone:%s\n", profile.Id.Hex(), profile.Phone)
	fmt.Println(userId)

	if userProfile, err := getUserProfile(cws, userId); err == nil {
		dbgMessage("%s: getprofile success\n", profile.Id.Hex())
		c.JSON(http.StatusOK, userProfile)

	} else {
		RespErrorMSG := make(map[string]string)
		RespErrorMSG["Message"] = "get profile fail"
		c.JSON(http.StatusBadRequest, RespErrorMSG)
	}
}

func (cws *careWorkerServer) profile(c *gin.Context) {
	profile := new(user_profile)
	// Obtain the POSTed JSON username and password values
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"payload": err.Error()})
		return
	}
	dbgMessage("profile JSON name:%s, phone:%s\n", profile.Id.Hex(), profile.Phone)
	fmt.Println(profile)

	if err := updateUserProfile(cws, profile.UserId, profile); err == nil {
		dbgMessage("%s: register success\n", profile.Id.Hex())
		c.JSON(http.StatusOK, "Success")

	} else {
		RespErrorMSG := make(map[string]string)
		RespErrorMSG["Message"] = "update profile fail"
		c.JSON(http.StatusBadRequest, RespErrorMSG)
	}
}

func (cws *careWorkerServer) performLogin(c *gin.Context) {
	loginUserAccount := new(user_account)

	// Obtain the POSTed JSON username and password values
	if err := c.ShouldBindJSON(&loginUserAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"payload": err.Error()})
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
		dbgMessage("set %s to session", user.Username)
		err := session.Save()
		if err != nil {
			dbgMessage("user session svae failed")
		}
		RespUser := make(map[string]string)
		RespUser["Username"] = user.Username
		RespUser["UserId"] = user.Id.Hex()
		dbgMessage("RespUser[Username]:%s", RespUser["Username"])
		dbgMessage("RespUser[UserId]:%s", RespUser["UserId"])
		c.JSON(http.StatusOK, RespUser)
	} else {
		RespErrorMSG := make(map[string]string)
		RespErrorMSG["Message"] = "login fail, Please chech account and password"
		c.JSON(http.StatusUnauthorized, RespErrorMSG)
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
		c.JSON(http.StatusBadRequest, gin.H{"payload": err.Error()})
		return
	}

	if _, err := registerNewUser(cws, newUserAccount); err == nil {
		dbgMessage("%s: register success\n", newUserAccount.Email)
		c.JSON(http.StatusOK, "Success")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		RespErrorMSG := make(map[string]string)
		RespErrorMSG["Message"] = err.Error()
		c.JSON(http.StatusBadRequest, RespErrorMSG)
	}
}

func (cws *careWorkerServer) registerSalt(c *gin.Context) {
	// Obtain the POSTed JSON username and password values
	saltUser := new(user_account)
	if err := c.ShouldBindJSON(&saltUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"payload": err.Error()})
		return
	}
	dbgMessage("register JSON email:%s\n", saltUser.Email)

	if queryUser, err := isUserSaltAvailable(cws, saltUser.Email); err == true {
		Salt := make(map[string]string)
		Salt["salt"] = queryUser.Salt
		Salt["username"] = queryUser.Username
		c.JSON(http.StatusOK, Salt)
	}
}

func (cws *careWorkerServer) forgotPassword(c *gin.Context) {
	userEmail := c.Param("email")

	dbgMessage("forgot email :%s\n", userEmail)
	// save reset hashcode in user database
	resetCode, err := setResetCode(cws, userEmail)

	if err == false {
		c.JSON(http.StatusOK, "Fail")
	}

	// send mail
	sendPasswordResetMail(userEmail, resetCode)

	c.JSON(http.StatusOK, "Success")
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
		c.JSON(http.StatusOK, "Verify Failed")
	} else {
		// clear resetCode
		clearResetCode(cws, userEmail)
		// update password
		resetPassword(cws, userEmail, userNewPassword)
		c.JSON(http.StatusOK, "Success")
	}
}
