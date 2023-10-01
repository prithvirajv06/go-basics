package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
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
var secretKey = []byte("your-secret-key")

func JWTMiddleware(next httprouter.Handle) httprouter.Handle {
	secretKey = []byte(os.Getenv("jwt_key"))
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		tokenString := r.Header.Get("Authorization")
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
		next(w, r, params)
	}
}

func createTokenForUser(userid string) (string, error) {
	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims (payload) for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "GBTT_USER" // Subject
	claims["userid"] = userid
	claims["iat"] = time.Now().Unix()                    // Issued At Time
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Expiration Time (1 hour from now)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUid() string {
	var date string = time.Now().Local().String()
	uid := strings.ReplaceAll(date, " ", "")
	reg := regexp.MustCompile(`[^a-zA-Z0-9\\s]+`)
	// Replace all matched characters with an empty string
	cleaned := reg.ReplaceAllString(uid, "")
	return cleaned
}
