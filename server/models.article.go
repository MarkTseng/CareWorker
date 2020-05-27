// models.article.go

package main

import (
	"errors"
	"github.com/night-codes/mgo-ai"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	//"log"
	//"fmt"
	//"strconv"
)

// Article model
type article struct {
	Id bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	//Id        uint64 `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string    `json:"title" form:"title" binding:"required" bson:"title"`
	Body      string    `json:"body" form:"body" binding:"required" bson:"body"`
	Location  string    `json:"location" form:"location" binding:"required" bson:"location"`
	Salary    string    `json:"salary" form:"salary" binding:"required" bson:"salary"`
	Author    string    `json:"author" form:"author" binding:"required" bson:"author"`
	CreatedOn time.Time `json:"created_on" bson:"created_on"`
	UpdatedOn time.Time `json:"updated_on" bson:"updated_on"`
	State     int       `json:"state"`
}

func getAllArticles(cws *careWorkerServer) []article {
	var results []article
	cws.collection["articles"].Find(nil).Sort("-timestamp").Limit(10).All(&results)

	return results
}

func getArticleByID(cws *careWorkerServer, id bson.ObjectId) ([]article, error) {
	result := article{}

	cws.collection["articles"].Find(bson.M{"_id": id}).One(&result)
	if result.Id == id {
		resultArray := []article{result}
		return resultArray, nil
	}

	return nil, errors.New("Article not found")
}

func createNewArticle(cws *careWorkerServer, title, content, location, salary, username string) (*article, error) {
	// update article counter
	ai.Connect(cws.collection["counters"])
	ai.Next("articles")

	aId := bson.NewObjectId()
	createdOnTime := time.Now()
	a := article{Title: title, Body: content, Id: aId, Author: username, Location: location, Salary: salary, CreatedOn: createdOnTime, UpdatedOn: createdOnTime}
	err := cws.collection["articles"].Insert(bson.M{"_id": aId, "title": title, "body": content, "author": username, "location": location, "salary": salary, "created_on": createdOnTime, "updated_on": createdOnTime})
	if err != nil {
		panic(err)
	}

	return &a, nil
}

func deleteOldArticle(cws *careWorkerServer, id, username string) error {
	//article_Id, err := strconv.Atoi(id)
	article_Id := bson.ObjectIdHex(id)
	err := cws.collection["articles"].Remove(bson.M{"_id": article_Id})
	if err != nil {
		dbgMessage("remove fail %v\n", err)
	}
	return err
}

func updateOldArticle(cws *careWorkerServer, id, title, content, username string) (*article, error) {
	//article_Id, err := strconv.Atoi(id)
	article_Id := bson.ObjectIdHex(id)
	ietmSelector := bson.M{"_id": article_Id}
	change := bson.M{"$set": bson.M{"title": title, "body": content, "author": username}}
	err := cws.collection["articles"].Update(ietmSelector, change)

	if err != nil {
		dbgMessage("update fail %v\n", err)
	}

	//a := article{Title: title, Body: content, Id: uint64(article_Id), Author: username}
	a := article{Title: title, Body: content, Id: article_Id, Author: username}
	return &a, nil
}
