package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//ConnecToDB: Connection to MongoDB database
func ConnecToDB() (*mgo.Collection, *mgo.Session) {
	sess, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	sess.SetMode(mgo.Monotonic, true)
	c := sess.DB("example").C("resources")
	return c, sess
}

//UpdateDBViaId Update data by using resource.Id
func UpdateDBViaId(resource res) error {
	c, sess := ConnecToDB()
	err := c.Update(bson.M{"X": resource.X, "Y": resource.Y, "Priority": resource.Priority},
		bson.M{"X": resource.X, "Y": resource.Y, "Priority": resource.Priority})
	defer sess.Close()
	return err
}

//InsertData Insert data to MongoDB database
func InsertData(resource res) error {
	c, sess := ConnecToDB()
	err := c.Insert(&resource)
	defer sess.Close()
	return err
}

//DeleteData Delete data from MongoDB database
func DeleteData(resource res) error {
	c, sess := ConnecToDB()
	err := c.Remove(bson.M{"ID": resource.ID})
	defer sess.Close()
	return err
}

//FindData Find one data in MongoDB database
func FindData(Id int) (result res, err error) {
	c, sess := ConnecToDB()
	err = c.Find(bson.M{"ID": Id}).One(&result)
	defer sess.Close()
	return result, err
}
