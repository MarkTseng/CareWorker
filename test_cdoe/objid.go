package main

import (
	"fmt"
	"github.com/pborman/uuid"
	"labix.org/v2/mgo/bson"
)

func main() {
	var newId bson.ObjectId

	newId = bson.NewObjectId()
	fmt.Println(newId)
	newId = bson.ObjectIdHex("5df845e46e95522363000001")
	fmt.Println(newId)
	fmt.Println(uuid.New())
}
