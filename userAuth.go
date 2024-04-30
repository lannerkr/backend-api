package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	// jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("!DE*DD&@#JUsecret"), nil)

func authHandler(w http.ResponseWriter, r *http.Request) {
	var u AuthUsers
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	dbconnCheck(dbconn)
	coll := dbconn.Database("usertableDB").Collection("authUsers")
	var dbuser AuthUsers

	filter := bson.D{{Key: "username", Value: u.Username}}
	user := coll.FindOne(context.TODO(), filter)
	err = user.Decode(&dbuser)
	if err != nil {
		log.Println("getAuthUserDB Find Error", err)
		return
	}

	// hashedpw := hashresult(u.Password)

	// if hashedpw != dbuser.Password {
	// 	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	// 	return
	// }

	acl := configuration.IPaccess
	rip := strings.Split(r.RemoteAddr, ":")[0]
	log.Printf("login user :%v , remoteAddr :%v", u.Username, rip)
	if slices.Contains(acl, rip) {
		log.Println("rip ", rip, " is accepted")
	} else {
		result := fmt.Sprint("rip ", rip, " is not in the IPaccess List")
		log.Println(result)
		http.Error(w, result, http.StatusForbidden)
		return
	}

	if !CheckPasswordHash(u.Password, dbuser.Password) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	claims := map[string]interface{}{"authorized": true, "user_id": dbuser.UserID, "isadmin": dbuser.IsAdmin}
	jwtauth.SetExpiryIn(claims, time.Hour)
	// "exp": time.Now().Add(time.Minute * 15).Unix()
	_, token, err := tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	responseUser := AuthToken{
		UserID:   dbuser.UserID,
		Username: dbuser.Username,
		IsAdmin:  dbuser.IsAdmin,
		JWTtoken: token,
		Exp:      time.Now().Add(time.Minute * 15).Unix(),
	}

	render.JSON(w, r, responseUser)
}

// func updateAuthToken(w http.ResponseWriter, r *http.Request, jwtToken jwt.Token) {
// 	v := AuthToken{
// 		JWTtoken: fmt.Sprint(jwtToken),
// 	}
// 	render.JSON(w, r, v)
// }

type AuthToken struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	JWTtoken string `json:"jwtToken"`
	Exp      int64  `json:"exp"`
}

// func createToken(userId int, isAdmin bool) (string, error) {
// 	var err error
// 	//Creating Access Token
// 	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
// 	atClaims := jwt.MapClaims{}
// 	atClaims["authorized"] = true
// 	atClaims["user_id"] = userId
// 	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
// 	atClaims["isadmin"] = isAdmin
// 	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
// 	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
// 	if err != nil {
// 		return "", err
// 	}
// 	return token, nil
// }

type AuthUsers struct {
	//	Id        bson.ObjectId `json:"_id,omitempty" bson:"_id"`
	UserID   int    `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	// FirstName string `json:"firstName" bson:"firstName"`
	// LastName  string `json:"lastName" bson:"lastName"`
	IsAdmin  bool   `json:"isAdmin" bson:"isAdmin"`
	JWTtoken string `json:"jwttoken"`
}

// var user = AuthUsers{
// 	UserID:    1,
// 	Username:  "username",
// 	Password:  "password",
// 	FirstName: "49123454322",
// 	LastName:  "asdgodk",
// }

func getUsersHandler(w http.ResponseWriter, r *http.Request) {

	userid := chi.URLParam(r, "userid")

	//var rescode []AuthUsers
	rescode := getAuthUserDB(userid)
	//fmt.Println(rescode)

	respondwithJSON(w, 200, rescode)
	//render.JSON(w, r, rescode)
}

func putUserHandler(w http.ResponseWriter, r *http.Request) {

	userid := chi.URLParam(r, "userid")

	var u AuthUsers
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	//hashedpw := hashresult(u.Password)
	hashedpw, err := HashPassword(u.Password)
	if err != nil {
		log.Println("userAuth password hashing error :" + err.Error())
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	if userid == "" {
		u.Password = hashedpw
		err = putAuthUserDB(u)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	} else {
		u.UserID, _ = strconv.Atoi(userid) //strconv.ParseUint(userid, 10, 64)
		u.Password = hashedpw
		err = updateAuthUserDB(u)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}

	respondwithJSON(w, http.StatusOK, http.StatusOK)
	//render.JSON(w, r, rescode)
}

func delUserHandler(w http.ResponseWriter, r *http.Request) {

	userid := chi.URLParam(r, "userid")

	//var rescode []AuthUsers
	rescode := delAuthUserDB(userid)
	//fmt.Println(rescode)

	respondwithJSON(w, 200, rescode)
	//render.JSON(w, r, rescode)
}

func getAuthUserDB(userid string) []AuthUsers {
	dbconnCheck(dbconn)
	coll := dbconn.Database("usertableDB").Collection("authUsers")

	var authuser []AuthUsers
	var filter primitive.D
	// u, err := coll.Find(context.TODO(), filter)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil
	// }

	if userid == "all" {
		filter = bson.D{}
	} else {
		uid, _ := strconv.Atoi(userid)
		filter = bson.D{{Key: "id", Value: uid}}
	}

	u, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.Println("getAuthUserDB Find Error", err)
		return nil
	}

	if err = u.All(context.TODO(), &authuser); err != nil {
		log.Println("getAuthUserDB Bind Error", err)
		return nil
	}

	//log.Println(users)
	return authuser
}

func putAuthUserDB(authuser AuthUsers) error {
	dbconnCheck(dbconn)
	coll := dbconn.Database("usertableDB").Collection("authUsers")

	//var authuser []AuthUsers
	groupStage := bson.D{{Key: "$group", Value: bson.M{"_id": "null", "max": bson.D{{Key: "$max", Value: "$id"}}}}}

	// {"$group", bson.D{
	// 	{"_id", "$category"},
	// 	{"average_price", bson.D{{"$avg", "$price"}}},
	// 	{"type_total", bson.D{{"$sum", 1}}},
	// }}}
	// pass the pipeline to the Aggregate() method
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{groupStage})
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	maxid, _ := strconv.Atoi(fmt.Sprint(results[0]["max"]))
	log.Printf("aggregate result : %v,%v\n", results[0]["max"], maxid)
	authuser.UserID = maxid + 1

	if _, err := coll.InsertOne(context.TODO(), authuser); err != nil {
		log.Println(err)
		return err
	}

	//log.Println(users)
	return nil
}

func updateAuthUserDB(authuser AuthUsers) error {
	dbconnCheck(dbconn)
	coll := dbconn.Database("usertableDB").Collection("authUsers")

	//var authuser []AuthUsers

	filter := bson.D{{Key: "id", Value: authuser.UserID}, {Key: "username", Value: authuser.Username}}
	update := bson.D{{Key: "$set", Value: bson.M{"password": authuser.Password}}}
	if _, err := coll.UpdateOne(context.TODO(), filter, update); err != nil {
		log.Println(err)
		return err
	}

	//log.Println(users)
	return nil
}

func delAuthUserDB(userid string) error {
	dbconnCheck(dbconn)
	coll := dbconn.Database("usertableDB").Collection("authUsers")

	uid, _ := strconv.Atoi(userid)

	filter := bson.D{{Key: "id", Value: uid}}

	if _, err := coll.DeleteOne(context.TODO(), filter); err != nil {
		log.Println(err)
		return err
	}

	//log.Println(users)
	return nil
}
