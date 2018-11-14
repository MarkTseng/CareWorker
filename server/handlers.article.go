// handlers.article.go

package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	//"log"
	"net/http"
	"strconv"
)

func showIndexPage(c *gin.Context) {
	articles := getAllArticles()

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title":   "Care Worker",
		"payload": articles}, "index.html")
}

func showArticleCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Create New Article"}, "create-article.html")
}

func getArticle(c *gin.Context) {
	// Check if the article ID is valid
	//log.Print(c.Param("article_id"))
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		// Check if the article exists
		if article, err := getArticleByID(uint64(articleID)); err == nil {
			// Call the render function with the title, article and the name of the
			// template
			render(c, gin.H{
				"title":   article.Title,
				"payload": article}, "article.html")

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func createArticle(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")

	// get username in session
	session := sessions.Default(c)
	username := session.Get("username")
	//log.Printf("createArticle username %s\n", username)

	if username != nil {
		if a, err := createNewArticle(title, content, username.(string)); err == nil {
			// If the article is created successfully, show success message
			render(c, gin.H{
				"title":   "Submission Successful",
				"payload": a}, "submission-successful.html")
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}

}

func deleteArticle(c *gin.Context) {
	id := c.Param("id")
	//log.Printf("deleteArticle id: %s\n", id)

	// get username in session
	session := sessions.Default(c)
	username := session.Get("username")
	//log.Printf("deleteArticle username %s\n", username)

	if username != nil {
		if err := deleteOldArticle(id, username.(string)); err == nil {
			// If the article is delete successfully, show success message
			render(c, gin.H{
				"title": "Submission Successful",
			}, "submission-delete-successful.html")
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}

}

func updateArticle(c *gin.Context) {
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
		if a, err := updateOldArticle(id, title, content, username.(string)); err == nil {
			// If the article is created successfully, show success message
			render(c, gin.H{
				"title":   "Submission Successful",
				"payload": a}, "submission-successful.html")
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
