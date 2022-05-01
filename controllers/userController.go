package controllers

import (
	"PersonalTaskManagement/database"
	"PersonalTaskManagement/entities"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var User entities.User
	json.NewDecoder(r.Body).Decode(&User)
	if !isEmailValid(User.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Email is invalid!")
		return
	}
	result := database.Instance.Create(&User)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result.Error.Error())
		return
	}
	json.NewEncoder(w).Encode("user_id:" + strconv.Itoa(User.ID) + " created Successfully!")

}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	UserId := mux.Vars(r)["id"]
	if checkIfUserExists(UserId) == false {
		json.NewEncoder(w).Encode("User Not Found!")
		return
	}
	var User entities.User
	database.Instance.First(&User, UserId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(User)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var Users []entities.User
	database.Instance.Find(&Users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	UserId := mux.Vars(r)["id"]
	if checkIfUserExists(UserId) == false {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("User Not Found!")
		return
	}
	var User entities.User
	database.Instance.Delete(&User, UserId)
	json.NewEncoder(w).Encode("User Deleted Successfully!")
}

func checkIfUserExists(UserId string) bool {
	var User entities.User
	database.Instance.First(&User, UserId)
	if User.ID == 0 {
		return false
	}
	return true
}
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
