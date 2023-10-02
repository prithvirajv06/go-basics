package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

func loggerStartup() *os.File {
	f, err := os.OpenFile("application", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	return f
}

func getEnvForStarup() (string, string) {
	godotenv.Load()
	dbConStr := os.Getenv("mongo_db_uri")
	port := os.Getenv("application_port")
	return dbConStr, port

}

func decodeRequestBody(r *http.Request, targetVar interface{}) {
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(targetVar)
}

// Define your secret key
var secretKey = []byte("")

func JWTMiddleware(next httprouter.Handle) httprouter.Handle {
	secretKey = []byte(os.Getenv("jwt_key"))
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Token is not valid", http.StatusUnauthorized)
			return
		}
		jwtCLaim, _ := parseToken(tokenString)
		w.Header().Add("userid", jwtCLaim.Sub)
		next(w, r, params)
	}
}

func createTokenForUser(userid string) (string, error) {
	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims (payload) for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userid // Subject
	claims["name"] = userid
	claims["iat"] = time.Now().Unix()                    // Issued At Time
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Expiration Time (1 hour from now)
	secretKey = []byte(os.Getenv("jwt_key"))

	// Sign the token with a secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseToken(tokenString string) (*JwtClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("jwt_key")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	if claims, ok := token.Claims.(*JwtClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("Failed to extract claims")
}

func getUid() string {
	var date string = time.Now().Local().String()
	uid := strings.ReplaceAll(date, " ", "")
	reg := regexp.MustCompile(`[^a-zA-Z0-9\\s]+`)
	// Replace all matched characters with an empty string
	cleaned := reg.ReplaceAllString(uid, "")
	return cleaned
}

func createBSONWithNonEmptyFields(data interface{}) (bson.M, error) {
	// Marshal the struct into a BSON map
	bsonData, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	// Unmarshal the BSON map into a BSON M document
	var bsonMap bson.M
	err = bson.Unmarshal(bsonData, &bsonMap)
	if err != nil {
		return nil, err
	}
	// Remove empty fields from the BSON M document
	for key, value := range bsonMap {
		if value == nil || value == "" {
			delete(bsonMap, key)
		}
	}
	return bsonMap, nil
}
