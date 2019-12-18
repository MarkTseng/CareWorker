// middleware.auth.go

package main

import (
	//	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// if the user is not logged in
func (cws *careWorkerServer) ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbgMessage("ensureLoggedIn")
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// if the user is already logged in
func (cws *careWorkerServer) ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbgMessage("ensureNotLoggedIn")
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if loggedIn {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// This middleware sets whether the user is logged in or not
func (cws *careWorkerServer) setUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			c.Set("is_logged_in", true)
			dbgMessage("setUserStatus is_logged_in true")
		} else {
			c.Set("is_logged_in", false)
			dbgMessage("setUserStatus is_logged_in false")
		}
	}
}
