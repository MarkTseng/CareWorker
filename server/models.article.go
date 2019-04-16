// models.article.go

package main

import (
	"errors"
	"github.com/night-codes/mgo-ai"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"strconv"
)

// Article model
type article struct {
	//Id        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Id        uint64 `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string `json:"title" form:"title" binding:"required" bson:"title"`
	Body      string `json:"body" form:"body" binding:"required" bson:"body"`
	Location  string `json:"location" form:"location" binding:"required" bson:"location"`
	Salary    string `json:"salary" form:"salary" binding:"required" bson:"salary"`
	Author    string `json:"author" form:"author" binding:"required" bson:"author"`
	CreatedOn int64  `json:"created_on" bson:"created_on"`
	UpdatedOn int64  `json:"updated_on" bson:"updated_on"`
}

// Return a list of all the articles
func getAllArticles(cws *careWorkerServer) []article {
	var results []article
	//cws.articles.Find(nil).Sort("-timestamp").All(&results)
	cws.collection["articles"].Find(nil).Sort("-timestamp").All(&results)

	return results
	//return articleList
}

// Fetch an article based on the Id supplied
func getArticleByID(cws *careWorkerServer, id uint64) (*article, error) {
	result := article{}
	//cws.articles.Find(bson.M{"_id": id}).One(&result)
	cws.collection["articles"].Find(bson.M{"_id": id}).One(&result)
	if result.Id == id {
		return &result, nil
	}

	return nil, errors.New("Article not found")
}

// Create a new article with the title and content provided
func createNewArticle(cws *careWorkerServer, title, content, location, salary, username string) (*article, error) {
	ai.Connect(cws.collection["counters"])
	aId := ai.Next("articles")
	a := article{Title: title, Body: content, Id: aId, Author: username, Location: location, Salary: salary}
	//err := cws.articles.Insert(bson.M{"_id": aId, "title": title, "body": content, "author": username})
	err := cws.collection["articles"].Insert(bson.M{"_id": aId, "title": title, "body": content, "author": username, "location": location, "salary": salary})
	if err != nil {
		panic(err)
	}

	return &a, nil
}

// Delete a old article with the title and content provided
func deleteOldArticle(cws *careWorkerServer, id, username string) error {
	article_Id, err := strconv.Atoi(id)
	//err = cws.articles.Remove(bson.M{"_id": article_Id})
	err = cws.collection["articles"].Remove(bson.M{"_id": article_Id})
	if err != nil {
		dbgMessage("remove fail %v\n", err)
	}
	return err
}

// Create a new article with the title and content provided
func updateOldArticle(cws *careWorkerServer, id, title, content, username string) (*article, error) {
	article_Id, err := strconv.Atoi(id)
	ietmSelector := bson.M{"_id": article_Id}
	change := bson.M{"$set": bson.M{"title": title, "body": content, "author": username}}
	//err = cws.articles.Update(ietmSelector, change)
	err = cws.collection["articles"].Update(ietmSelector, change)

	if err != nil {
		dbgMessage("update fail %v\n", err)
	}

	a := article{Title: title, Body: content, Id: uint64(article_Id), Author: username}
	return &a, nil
}
