package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"rudrprasad.com/backend/database"
	"rudrprasad.com/backend/logs"
)

type Routes struct {
	DB *sql.DB
	LOG *logs.Logger
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var jwtSecret = []byte("your_secret_key")

func GenerateJWT(email string) (string, error) {
	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	return token.SignedString(jwtSecret)
}

func (routes *Routes) PostLoginUser(w http.ResponseWriter, r *http.Request){
	var user database.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data := Message{Data: "invalid json"}
		routes.LOG.Error("400 bad request; invalid json")
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	if user.Email == "" || user.Password == ""{
		data := Message{Data: "missing information"}
		routes.LOG.Error("400 bad request; missing creds")
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	checkIfUserExist := database.GetUserByEmail(routes.DB,user.Email)
	if checkIfUserExist == nil {
		data := Message{Data: "user doesnt exist"}
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(checkIfUserExist.Password), []byte(user.Password))
	if err != nil{
		data := Message{Data: "invalid creds"}
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user.Email)
	if err != nil {
		data := Message{Data: "internal server error"}
		routes.LOG.Error("500 internal server error; failed to generate token")
		sendJSONResponse(w, http.StatusInternalServerError, data)
		return
	}

	response := map[string]string{
		"token": token,
	}
	sendJSONResponse(w, http.StatusOK, response)

}


func (routes *Routes) PostRegisterUser(w http.ResponseWriter, r *http.Request){
	var user database.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data := Message{Data: "invalid json"}
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	if user.Email == "" || user.Name == "" || user.Password == ""{
		data := Message{Data: "missing information"}
		routes.LOG.Error("400 bad request; missing creds")
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	checkIfUserExist := database.GetUserByEmail(routes.DB,user.Email)
	if checkIfUserExist != nil {
		data := Message{Data: "user exists"}
		routes.LOG.Error("400 bad request; user exists")
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
        data := Message{Data: "missing information"}
		routes.LOG.Error("500 internal error; unable to hashpassword")
		sendJSONResponse(w, http.StatusInternalServerError, data)
		return
    }

	user.Password = string(hashedPassword)
	newUser, newUserErr := database.CreateNewUser(routes.DB, &user)
	if newUserErr != nil{
		data := Message{Data: "user not created"}
		sendJSONResponse(w, http.StatusInternalServerError, data)
		return
	}

	sendJSONResponse(w, http.StatusOK, newUser)
	return
}