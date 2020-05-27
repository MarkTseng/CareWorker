// handlers.article.go

package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func (cws *careWorkerServer) showIndexPage(c *gin.Context) {
	articles := getAllArticles(cws)
	c.SecureJSON(http.StatusOK, articles)
}

func (cws *careWorkerServer) getArticle(c *gin.Context) {
	// Check if the article ID is valid
	if bson.IsObjectIdHex(c.Param("article_id")) == true {
		articleObjID := bson.ObjectIdHex(c.Param("article_id"))
		// Check if the article exists
		if article, err := getArticleByID(cws, articleObjID); err == nil {
			c.SecureJSON(http.StatusOK, article)

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (cws *careWorkerServer) createArticle(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")
	location := c.PostForm("location")
	salary := c.PostForm("salary")

	// Obtain the POSTed JSON username and password values
	if title == "" {
		objA := new(article)
		if err := c.ShouldBindJSON(&objA); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		title = objA.Title
		content = objA.Body
		location = objA.Location
		salary = objA.Salary
	}

	// get username in session
	session := sessions.Default(c)
	username := session.Get("username")

	if username != nil {
		if a, err := createNewArticle(cws, title, content, location, salary, username.(string)); err == nil {
			// If the article is created successfully, show success message
			/*
				render(c, gin.H{
					"title":   "Submission Successful",
					"payload": a.Id}, "submission-successful.html",
					http.StatusOK)
			*/
			c.SecureJSON(http.StatusOK, a.Id)
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}

}

func (cws *careWorkerServer) deleteArticle(c *gin.Context) {
	id := c.Param("id")
	dbgMessage("deleteArticle id: %s\n", id)

	// get username in session
	session := sessions.Default(c)
	username := session.Get("username")

	if username != nil {
		if err := deleteOldArticle(cws, id, username.(string)); err == nil {
			// If the article is delete successfully, show success message
			RespErrorMSG := make(map[string]string)
			RespErrorMSG["Message"] = "delete successful"
			c.SecureJSON(http.StatusOK, RespErrorMSG)
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}

}

func (cws *careWorkerServer) updateArticle(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")
	id := c.PostForm("id")

	// get username in session
	session := sessions.Default(c)
	username := session.Get("username")

	if username != nil {
		if a, err := updateOldArticle(cws, id, title, content, username.(string)); err == nil {
			// If the article is created successfully, show success message
			c.SecureJSON(http.StatusOK, a)
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
