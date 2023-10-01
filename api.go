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
	router.POST("/create-user", JWTMiddleware(CreateUser))
	return router
}

func WelcomeUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Welcoem To Go Buddies !"))
}

func LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var user *GBUser = new(GBUser)
	decodeRequestBody(r, &user)
	token, _ := createTokenForUser(user.GOBID)
	var commonRes *GBCommongResponse = new(GBCommongResponse)
	commonRes.Token = string(token)
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
