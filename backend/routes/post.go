package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"rudrprasad.com/backend/database"
	"rudrprasad.com/backend/logs"
)

type Routes struct {
	DB *sql.DB
	LOG *logs.Logger
}


func (routes *Routes) PostRegisterUser(w http.ResponseWriter, r *http.Request){
	var user database.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		data := Message{Data: "invalid json"}
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	checkIfUserExist := database.GetUserByEmail(routes.DB,user.Email)
	if checkIfUserExist != nil {
		data := Message{Data: "user exists"}
		sendJSONResponse(w, http.StatusBadRequest, data)
		return
	}

	newUser, newUserErr := database.CreateNewUser(routes.DB, &user)
	if newUserErr != nil{
		data := Message{Data: "user not created"}
		sendJSONResponse(w, http.StatusInternalServerError, data)
		return
	}

	sendJSONResponse(w, http.StatusOK, newUser)
	return
}