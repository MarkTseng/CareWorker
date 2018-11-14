// models.article.go

package main

import (
	"errors"
	"github.com/night-codes/mgo-ai"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
)

/*
type article struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Body string `json:"content"`
}
*/

// Article model
type article struct {
	//Id        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Id        uint64 `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string `json:"title" form:"title" binding:"required" bson:"title"`
	Body      string `json:"body" form:"body" binding:"required" bson:"body"`
	Author    string `json:"author" form:"author" binding:"required" bson:"author"`
	CreatedOn int64  `json:"created_on" bson:"created_on"`
	UpdatedOn int64  `json:"updated_on" bson:"updated_on"`
}

// For this demo, we're storing the article list in memory
// In a real application, this list will most likely be fetched
// from a database or from static files
/*
var articleList = []article{
	article{Id: bson.ObjectId("1"), Title: "Article 1", Body: "Article 1 body"},
	article{Id: bson.ObjectId("2"), Title: "Article 2", Body: "Article 2 body"},
}
*/
// Return a list of all the articles
func getAllArticles() []article {
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("articles")
	var results []article
	err = c.Find(nil).Sort("-timestamp").All(&results)

	return results
	//return articleList
}

// Fetch an article based on the Id supplied
func getArticleByID(id uint64) (*article, error) {
	/*
		for _, a := range articleList {
			if a.Id == bson.ObjectId(id) {
				return &a, nil
			}
		}
	*/
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("articles")
	result := article{}
	err = c.Find(bson.M{"_id": id}).One(&result)
	//log.Print(result)
	if result.Id == id {
		return &result, nil
	}

	return nil, errors.New("Article not found")
}

// Create a new article with the title and content provided
func createNewArticle(title, content, username string) (*article, error) {
	// Set the Id of a new article to one more than the number of articles
	//a := article{Title: title, Body: content}

	// Add the article to the list of articles
	// articleList = append(articleList, a)

	//log.Printf("createNewArticle: %s\n", username)
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	// connect AutoIncrement to collection "counters"
	ai.Connect(session.DB("articles_demo_dev").C("counters"))
	c := session.DB("articles_demo_dev").C("articles")

	aId := ai.Next("articles")
	a := article{Title: title, Body: content, Id: aId, Author: username}
	err = c.Insert(bson.M{"_id": aId, "title": title, "body": content, "author": username})
	//log.Print(aId)
	return &a, nil
}

// Delete a old article with the title and content provided
func deleteOldArticle(id, username string) error {
	// Set the Id of a new article to one more than the number of articles
	//a := article{Title: title, Body: content}

	// Add the article to the list of articles
	// articleList = append(articleList, a)

	log.Printf("deleteOldArticle")
	log.Printf("id:%s\n", id)
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("articles")

	//remove record
	article_Id, err := strconv.Atoi(id)
	err = c.Remove(bson.M{"_id": article_Id})
	if err != nil {
		log.Printf("remove fail %v\n", err)
	}
	//log.Print(aId)
	return err
}

// Create a new article with the title and content provided
func updateOldArticle(id, title, content, username string) (*article, error) {
	// Set the Id of a new article to one more than the number of articles
	//a := article{Title: title, Body: content}

	// Add the article to the list of articles
	// articleList = append(articleList, a)

	//log.Printf("createNewArticle: %s\n", username)
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("articles")

	// Update
	article_Id, err := strconv.Atoi(id)
	ietmSelector := bson.M{"_id": article_Id}
	change := bson.M{"$set": bson.M{"title": title, "body": content, "author": username}}
	err = c.Update(ietmSelector, change)

	if err != nil {
		log.Printf("update fail %v\n", err)
	}

	a := article{Title: title, Body: content, Id: uint64(article_Id), Author: username}
	//log.Print(aId)
	return &a, nil
}
