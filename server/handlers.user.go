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
)

func generateSessionToken() string {
	return uuid.New()
}

func (cws *careWorkerServer) performLogin(c *gin.Context) {
	loginUser := new(user)

	// Obtain the POSTed JSON username and password values
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"payload": err.Error()})
		return
	}
	dbgMessage("performLogin JSON email:%s, password:%s\n", loginUser.Email, loginUser.Password)

	// Check if the username/password combination is valid
	if user := isUserValid(cws, loginUser.Email, loginUser.Password); user != nil {
		// If the username/password is valid set the token in a cookie
		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

		// save username in session
		session := sessions.Default(c)
		session.Set("username", user.Username)
		session.Save()

		RespUser := make(map[string]string)
		RespUser["Username"] = user.Username
		RespUser["ID"] = user.Id.Hex()
		c.JSON(http.StatusOK, RespUser)
	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		render(c, gin.H{
			"payload": "Invalid credentials provided"},
			"login.html",
			http.StatusBadRequest)
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
	newUser := new(user)

	// Obtain the POSTed JSON username and password values
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"payload": err.Error()})
		return
	}

	if _, err := registerNewUser(cws, newUser); err == nil {
		dbgMessage("%s: register success\n", newUser.Email)
		c.JSON(http.StatusOK, "Success")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		render(c, gin.H{
			"payload": err.Error()},
			"register.html",
			http.StatusBadRequest)
	}
}

func (cws *careWorkerServer) registerSalt(c *gin.Context) {
	// Obtain the POSTed JSON username and password values
	saltUser := new(user)
	if err := c.ShouldBindJSON(&saltUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"payload": err.Error()})
		return
	}
	dbgMessage("register JSON email:%s\n", saltUser.Email)

	if queryUser, err := isUserSaleAvailable(cws, saltUser.Email); err == true {
		Salt := make(map[string]string)
		Salt["salt"] = queryUser.Salt
		Salt["username"] = queryUser.Username
		c.JSON(http.StatusOK, Salt)
	}
}
