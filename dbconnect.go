package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbconn *mongo.Client

func mongoConnect() *mongo.Client {
	uri := configuration.DBUri

	dbconn, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err)
	}

	return dbconn
}

func dbconnCheck(dbconn *mongo.Client) {
TOP:
	if err := dbconn.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println(err)
		log.Println("DBconn Fail and trying to reconnect")
		dbconn.Disconnect(context.TODO())
		time.Sleep(time.Second)
		dbconn = mongoConnect()
		goto TOP
	}
	//log.Println("DBconnection OK")
}

func getUsersHistory(realm string) (users []DBData) {

	dbconnCheck(dbconn)
	coll := dbconn.Database("usertableDB").Collection("user_history")

	//var users []DBData
	filter := bson.D{{Key: "realm", Value: realm}}
	u, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil
	}

	if err = u.All(context.TODO(), &users); err != nil {
		log.Println(err)
		return nil
	}

	//log.Println(users)
	return users
}

// func getUserTableDB(realm string) []UserData {

// 	dbconnCheck(dbconn)

// 	var users []UserData
// 	var realmColl string
// 	if realm == "EMP-GOTP" {
// 		realmColl = "userTable_" + "EMPGOTP"
// 	} else {
// 		realmColl = "userTable_" + realm
// 	}
// 	coll := dbconn.Database("ldapDB").Collection(realmColl)

// 	filter := bson.D{}
// 	u, err := coll.Find(context.TODO(), filter)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	if err = u.All(context.TODO(), &users); err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	return users
// }

// func putUserTableDB(realm string, users []UserData) error {

// 	dbconnCheck(dbconn)

// 	var realmColl string
// 	if realm == "EMP-GOTP" {
// 		realmColl = "userTable_" + "EMPGOTP"
// 	} else {
// 		realmColl = "userTable_" + realm
// 	}
// 	coll := dbconn.Database("ldapDB").Collection(realmColl)
// 	if err := coll.Drop(context.TODO()); err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	var update []interface{}
// 	for _, v := range users {
// 		v.Realm = realm
// 		update = append(update, v)
// 	}

// 	if _, err := coll.InsertMany(context.TODO(), update); err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	return nil
// }

// func getSingleUserTableDB(realm, user string) (userdb DBData) {

// 	dbconnCheck(dbconn)

// 	var usertable []UserData
// 	var realmColl string
// 	if realm == "EMP-GOTP" {
// 		realmColl = "userTable_" + "EMPGOTP"
// 	} else {
// 		realmColl = "userTable_" + realm
// 	}
// 	coll := dbconn.Database("ldapDB").Collection(realmColl)

// 	filter := bson.D{{Key: "user_name", Value: user}}
// 	u, err := coll.Find(context.TODO(), filter)
// 	if err != nil {
// 		log.Println(err)
// 		return userdb
// 	}

// 	if err = u.All(context.TODO(), &usertable); err != nil {
// 		log.Println(err)
// 		return userdb
// 	}

// 	return usertable[0].DBTable
// }

func putUserTableNewDB(realm string, users []UserData) error {

	dbconnCheck(dbconn)

	var realmColl string
	if realm == "EMP-GOTP" {
		realmColl = "userTable_" + "EMPGOTP"
	} else {
		realmColl = "userTable_" + realm
	}
	coll := dbconn.Database("usertableDB").Collection(realmColl)
	if err := coll.Drop(context.TODO()); err != nil {
		log.Println(err)
		return err
	}

	var update []interface{}
	for _, v := range users {
		v.Realm = realm
		update = append(update, v)
	}

	if _, err := coll.InsertMany(context.TODO(), update); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func getUserTableNewDB(realm string) []UserData {

	dbconnCheck(dbconn)

	var users []UserData
	var realmColl string
	if realm == "EMP-GOTP" {
		realmColl = "userTable_" + "EMPGOTP"
	} else {
		realmColl = "userTable_" + realm
	}
	coll := dbconn.Database("usertableDB").Collection(realmColl)

	filter := bson.D{}
	u, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil
	}

	if err = u.All(context.TODO(), &users); err != nil {
		log.Println(err)
		return nil
	}

	//log.Printf("getUserTableNewDB: \n%v\n", users[0])
	return users
}

func getSingleUserTableNewDB(realm, user string) (userdb DBData) {
	//log.Println("getSingleUserTableDB")

	dbconnCheck(dbconn)

	var usertable []UserData
	var realmColl string
	if realm == "EMP-GOTP" {
		realmColl = "userTable_" + "EMPGOTP"
	} else {
		realmColl = "userTable_" + realm
	}
	coll := dbconn.Database("usertableDB").Collection(realmColl)

	filter := bson.D{{Key: "user_name", Value: user}}
	u, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return userdb
	}

	if err = u.All(context.TODO(), &usertable); err != nil {
		log.Println(err)
		return userdb
	}

	return usertable[0].DBTable
}

// func getUsersOldHistory(realm string) (users []DBData) {

// 	dbconnCheck(dbconn)
// 	coll := dbconn.Database("ldapDB").Collection("user_history")

// 	//var users []DBData
// 	filter := bson.D{{Key: "realm", Value: realm}}
// 	u, err := coll.Find(context.TODO(), filter)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	if err = u.All(context.TODO(), &users); err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	//log.Println(users)
// 	return users
// }
