package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AllHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/", WelcomeUser)
	router.POST("/login", LoginUser)
	router.POST("/create-user", CreateUser)
	router.GET("/get-user", JWTMiddleware(GetUser))
	return router
}

func WelcomeUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Welcoem To Go Buddies !"))
}

func LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var user *GBUser = new(GBUser)
	decodeRequestBody(r, &user)
	_, err := findOne("user", &user)
	var commonRes *GBCommongResponse = new(GBCommongResponse)
	if err != nil {
		commonRes.Message = "User Not Found !"
		w.WriteHeader(http.StatusForbidden)
	} else {
		token, _ := createTokenForUser(user.GOBID)
		commonRes.Token = string(token)
		w.WriteHeader(http.StatusAccepted)
	}
	response, _ := json.Marshal(commonRes)
	w.Write(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var user *GBUser = new(GBUser)
	decodeRequestBody(r, &user)
	user.GOBID = getUid()
	insertOne("user", user)
	resonse, _ := json.Marshal(user)
	w.Write(resonse)
}

func GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var user *GBUser = new(GBUser)
	decodeRequestBody(r, &user)
	findOne("user", &user)
	resonse, _ := json.Marshal(*user)
	w.Write(resonse)
}
