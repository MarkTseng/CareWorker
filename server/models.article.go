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
	State     int    `json:"state"`
	Tag       string `json:"tag"`
}

func getAllArticles(cws *careWorkerServer) []article {
	var results []article
	cws.collection["articles"].Find(nil).Sort("-timestamp").Limit(10).All(&results)

	return results
}

func getArticleByID(cws *careWorkerServer, id uint64) ([]article, error) {
	result := article{}

	cws.collection["articles"].Find(bson.M{"_id": id}).One(&result)
	if result.Id == id {
		resultArray := []article{result}
		return resultArray, nil
	}

	return nil, errors.New("Article not found")
}

func createNewArticle(cws *careWorkerServer, title, content, location, salary, username string) (*article, error) {
	ai.Connect(cws.collection["counters"])
	aId := ai.Next("articles")
	a := article{Title: title, Body: content, Id: aId, Author: username, Location: location, Salary: salary}
	err := cws.collection["articles"].Insert(bson.M{"_id": aId, "title": title, "body": content, "author": username, "location": location, "salary": salary})
	if err != nil {
		panic(err)
	}

	return &a, nil
}

func deleteOldArticle(cws *careWorkerServer, id, username string) error {
	article_Id, err := strconv.Atoi(id)
	err = cws.collection["articles"].Remove(bson.M{"_id": article_Id})
	if err != nil {
		dbgMessage("remove fail %v\n", err)
	}
	return err
}

func updateOldArticle(cws *careWorkerServer, id, title, content, username string) (*article, error) {
	article_Id, err := strconv.Atoi(id)
	ietmSelector := bson.M{"_id": article_Id}
	change := bson.M{"$set": bson.M{"title": title, "body": content, "author": username}}
	err = cws.collection["articles"].Update(ietmSelector, change)

	if err != nil {
		dbgMessage("update fail %v\n", err)
	}

	a := article{Title: title, Body: content, Id: uint64(article_Id), Author: username}
	return &a, nil
}
