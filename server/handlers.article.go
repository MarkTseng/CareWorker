// handlers.article.go

package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (cws *careWorkerServer) showIndexPage(c *gin.Context) {
	articles := getAllArticles(cws)

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title":   "Care Worker",
		"payload": articles},
		"index.html",
		http.StatusOK)
}

func (cws *careWorkerServer) showArticleCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Create New Article"},
		"create-article.html",
		http.StatusOK)
}

func (cws *careWorkerServer) getArticle(c *gin.Context) {
	// Check if the article ID is valid
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		// Check if the article exists
		if article, err := getArticleByID(cws, uint64(articleID)); err == nil {
			// Call the render function with the title, article and the name of the template
			render(c, gin.H{
				"title":   article.Title,
				"payload": article}, "article.html",
				http.StatusOK)

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
	log.Printf("Title: %s\n", title)
	log.Printf("content: %s\n", content)
	log.Printf("location: %s\n", location)
	log.Printf("salary: %s\n", salary)

	// get username in session
	session := sessions.Default(c)
	username := session.Get("username")
	//log.Printf("createArticle username %s\n", username)

	if username != nil {
		if a, err := createNewArticle(cws, title, content, location, salary, username.(string)); err == nil {
			// If the article is created successfully, show success message
			render(c, gin.H{
				"title":   "Submission Successful",
				"payload": a.Id}, "submission-successful.html",
				http.StatusOK)
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
	//log.Printf("deleteArticle id: %s\n", id)

	// get username in session
	session := sessions.Default(c)
	username := session.Get("username")
	//log.Printf("deleteArticle username %s\n", username)

	if username != nil {
		if err := deleteOldArticle(cws, id, username.(string)); err == nil {
			// If the article is delete successfully, show success message
			render(c, gin.H{
				"title": "Submission Successful"},
				"submission-delete-successful.html",
				http.StatusOK)
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

	//log.Printf("updateArticle id %s\n", id)
	//log.Printf("updateArticle title %s\n", title)

	// get username in session
	session := sessions.Default(c)
	username := session.Get("username")
	//log.Printf("createArticle username %s\n", username)

	if username != nil {
		if a, err := updateOldArticle(cws, id, title, content, username.(string)); err == nil {
			// If the article is created successfully, show success message
			render(c, gin.H{
				"title":   "Submission Successful",
				"payload": a}, "submission-successful.html",
				http.StatusOK)
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
