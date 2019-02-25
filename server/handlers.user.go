// handlers.user.go

package main

import (
	//"math/rand"
	"net/http"
	//"strconv"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"log"
)

func generateSessionToken() string {
	// We're using a random 16 character string as the session token
	// This is NOT a secure way of generating session tokens
	// DO NOT USE THIS IN PRODUCTION
	//return strconv.FormatInt(rand.Int63(), 16)
	return uuid.New()
}

func (cws *careWorkerServer) showLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func (cws *careWorkerServer) performLogin(c *gin.Context) {
	// Obtain the POSTed username and password values
	email := c.PostForm("email")
	password := c.PostForm("password")
	loginUser := new(user)

	log.Printf("performLogin POST email:%s, password:%s\n", email, password)
	// Obtain the POSTed JSON username and password values
	if email == "" && password == "" {
		if err := c.ShouldBindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("performLogin JSON email:%s, password:%s\n", loginUser.Email, loginUser.Password)
		email = loginUser.Email
		password = loginUser.Password
	}

	// Check if the username/password combination is valid
	if user := isUserValid(cws, email, password); user != nil {
		// If the username/password is valid set the token in a cookie
		token := generateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

		// save username in session
		session := sessions.Default(c)
		session.Set("username", user.Username)
		session.Save()

		//log.Printf("username %s\n", username)
		RespUser := make(map[string]string)
		RespUser["Username"] = user.Username
		RespUser["ID"] = user.Id.Hex()
		c.JSON(http.StatusOK, RespUser)
	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
}

func (cws *careWorkerServer) logout(c *gin.Context) {
	log.Printf("logout\n")
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

func (cws *careWorkerServer) showRegistrationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	log.Println("showRegistrationPage")
	render(c, gin.H{
		"title": "Register"}, "register.html")
}

func (cws *careWorkerServer) register(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")
	salt := c.PostForm("salt")
	email := c.PostForm("email")
	newUser := new(user)

	log.Printf("register POST username:%s, password:%s, salt:%s\n", username, password, salt)
	// Obtain the POSTed JSON username and password values
	if username == "" && password == "" {
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println(newUser)
		username = newUser.Username
		password = newUser.Password
		salt = newUser.Salt
		email = newUser.Email
	}

	if _, err := registerNewUser(cws, username, password, salt, email, newUser); err == nil {
		/*
			// If the user is created, set the token in a cookie and log the user in
				token := generateSessionToken()
				c.SetCookie("token", token, 3600, "", "", false, true)
				c.Set("is_logged_in", true)

				// save username in session

				session := sessions.Default(c)
				session.Set("username", username)
				session.Save()
		*/
		log.Printf("%s: register success\n", username)
		c.JSON(http.StatusOK, "Success")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error()})

	}
}

func (cws *careWorkerServer) registerSalt(c *gin.Context) {
	// Obtain the POSTed username and password values
	email := c.PostForm("email")

	log.Printf("register POST username:%s\n")
	// Obtain the POSTed JSON username and password values
	if email == "" {
		saltUser := new(user)
		if err := c.ShouldBindJSON(&saltUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("register JSON email:%s\n", saltUser.Email)

		if queryUser, err := isUserSaleAvailable(cws, saltUser.Email); err == true {
			Salt := make(map[string]string)
			Salt["salt"] = queryUser.Salt
			Salt["username"] = queryUser.Username
			log.Println(queryUser)
			log.Println(err)
			c.JSON(http.StatusOK, Salt)
		}
	}
}
