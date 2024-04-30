package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PulseTable struct {
	User []UserData `json:"user"`
}
type APITable struct {
	User []UserData `json:"data"`
}

type UserData struct {
	Username  string `json:"username" bson:"user_name"`             // DB + PCS
	Enabled   string `json:"enabled" bson:"enabled"`                // DB + PCS
	Fullname  string `json:"fullname" bson:"fullname"`              // PCS
	PwdChange string `json:"change-password-at-signin" bson:"cpas"` // PCS

	Realm      string     `json:"realm,omitempty" bson:"realm"` // DB
	DBTable    DBData     `json:"dbtable" bson:"dbtable"`
	UserRecord Records    `json:"records"`
	UserLog    []UserLogs `json:"userlog"`
}

type DBData struct {
	Username     string    `json:"username" bson:"user_name"`
	Enabled      string    `json:"enabled" bson:"enabled"`
	LastLogin    time.Time `json:"lastLogin" bson:"last_login"`
	Realm        string    `json:"realm,omitempty" bson:"realm"`
	Days         int       `json:"days" bson:"days"`
	AccExpires   time.Time `json:"accountExpires" bson:"accountExpires"`
	StaticIP     string    `json:"Static_IP" bson:"Static_IP"`
	FramedIP     string    `json:"framedip" bson:"framedip"`
	LoginName    string    `json:"login_name" bson:"login_name"`
	LastString   string    `json:"laststring" bson:"laststring"`
	ExpireString string    `json:"expirestring" bson:"expirestring"`
}

// func updateDBtable(realm string, force bool) {
// 	location, _ := time.LoadLocation("Asia/Seoul")
// 	time.Local = location
// 	if checkDBupdateNeed(realm) || force {
// 		pcsusers := getUsersPCS(realm)
// 		dbusers := getUsersHistory(realm)
// 		var dbtable map[string]DBData = make(map[string]DBData)
// 		for _, v := range dbusers {
// 			if v.AccExpires.Year() < 2000 || v.AccExpires.Year() > 3000 {
// 				v.ExpireString = ""
// 			} else {
// 				v.ExpireString = v.AccExpires.Local().Format("2006-01-02T15:04")
// 			}
// 			v.LastString = v.LastLogin.Local().Format("2006-01-02T15:04")
// 			dbtable[v.Username] = v
// 		}
// 		for i, v := range pcsusers {
// 			pcsusers[i].DBTable = dbtable[v.Username]
// 		}
// 		putUserTableDB(realm, pcsusers)

// 		collu := dbconn.Database("ldapDB").Collection("updateTime")
// 		filter := bson.M{"realm": realm}
// 		newdate := bson.D{{Key: "$set", Value: bson.M{"updateTime": time.Now()}}}
// 		opts := options.Update().SetUpsert(true)
// 		if _, err := collu.UpdateOne(context.TODO(), filter, newdate, opts); err != nil {
// 			log.Println(err)
// 		}

// 		log.Println("usertable DB updated")
// 	}
// }

type Utime struct {
	UpdateTime time.Time `bson:"updateTime"`
}

// func checkDBupdateNeed(realm string) bool {

// 	collu := dbconn.Database("ldapDB").Collection("updateTime")
// 	var result Utime
// 	if err := collu.FindOne(context.TODO(), bson.D{{Key: "realm", Value: realm}}).Decode(&result); err != nil {
// 		log.Printf("check Find fail: %v\n", err)
// 		return true
// 	}

// 	t1 := result.UpdateTime.Add(time.Hour)

// 	return t1.Before(time.Now())
// }

type Response struct {
	ActiveUsers ActiveUserRecords `json:"active-users"`
}

type ActiveUserRecords struct {
	ActiveUserRecord UserRecords `json:"active-user-records"`
	TotalMatch       int         `json:"total-matched-record-number"`
	TotalRuturn      int         `json:"total-returned-record-number"`
	UserPermission   bool        `json:"user-login-permission"`
}

type UserRecords struct {
	UserRecord []Records `json:"active-user-record"`
}

type Records struct {
	Username      string `json:"active-user-name"`
	Realm         string `json:"authentication-realm"`
	ConnectIP     string `json:"network-connect-ip"`
	ClientVersion string `json:"pulse-client-version"`
	Role          string `json:"user-roles"`
	LoginTime     string `json:"user-sign-in-time"`
	SessionID     string `json:"session-id"`
}

func updateNewDBtable(realm string, force bool) {
	location, _ := time.LoadLocation("Asia/Seoul")
	time.Local = location
	if checkNewDBupdateNeed(realm) || force {
		pcsusers := getUsersPCS(realm)
		dbusers := getUsersHistory(realm)
		var dbtable map[string]DBData = make(map[string]DBData)
		for _, v := range dbusers {
			if v.AccExpires.Year() < 2000 || v.AccExpires.Year() > 3000 {
				v.ExpireString = ""
			} else {
				v.ExpireString = v.AccExpires.Local().Format("2006-01-02T15:04")
			}
			v.LastString = v.LastLogin.Local().Format("2006-01-02T15:04")
			dbtable[v.Username] = v
		}
		for i, v := range pcsusers {
			pcsusers[i].DBTable = dbtable[v.Username]
		}
		putUserTableNewDB(realm, pcsusers)

		collu := dbconn.Database("usertableDB").Collection("updateTime")
		filter := bson.M{"realm": realm}
		newdate := bson.D{{Key: "$set", Value: bson.M{"updateTime": time.Now()}}}
		opts := options.Update().SetUpsert(true)
		if _, err := collu.UpdateOne(context.TODO(), filter, newdate, opts); err != nil {
			log.Println(err)
		}

		log.Println("usertable DB updated")
	}
}

func checkNewDBupdateNeed(realm string) bool {

	collu := dbconn.Database("usertableDB").Collection("updateTime")
	var result Utime
	if err := collu.FindOne(context.TODO(), bson.D{{Key: "realm", Value: realm}}).Decode(&result); err != nil {
		log.Printf("check Find fail: %v\n", err)
		return true
	}

	t1 := result.UpdateTime.Add(time.Hour)

	return t1.Before(time.Now())
}
