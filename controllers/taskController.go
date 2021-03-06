package controllers

import (
	"PersonalTaskManagement/database"
	"PersonalTaskManagement/entities"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var Task entities.Task
	json.NewDecoder(r.Body).Decode(&Task)
	if Task.UserID == 0 || Task.EndTime.IsZero() || Task.StartTime.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("check input parameters!")
		return
	}

	if !checkIfUserExists(strconv.Itoa(Task.UserID)) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("UserID Not Found!")
		return
	}
	d1 := Task.EndTime.Sub(Task.StartTime).Seconds()
	if d1 < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Task End_Time is sooner than Start_Time!")
		return
	}
	if !checkIfUserHasTaskInPeriod(Task) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Task (Start-End Time) has overlap with other tasks!")
		return
	}
	result := database.Instance.Create(&Task)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result.Error.Error())
		return
	}

	json.NewEncoder(w).Encode("task ID:" + strconv.Itoa(int(Task.ID)) + " created Successfully!")

}

func GetTaskById(w http.ResponseWriter, r *http.Request) {
	TaskId := mux.Vars(r)["id"]
	if checkIfTaskExists(TaskId) == false {
		json.NewEncoder(w).Encode("Task Not Found!")
		return
	}
	var Task entities.Task
	database.Instance.First(&Task, TaskId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	var Tasks []entities.Task
	database.Instance.Find(&Tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Tasks)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	TaskId := mux.Vars(r)["id"]
	if checkIfTaskExists(TaskId) == false {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Task Not Found!")
		return
	}
	var Task entities.Task
	database.Instance.Delete(&Task, TaskId)
	json.NewEncoder(w).Encode("Task Deleted Successfully!")
}

func checkIfTaskExists(TaskId string) bool {
	var Task entities.Task
	database.Instance.First(&Task, TaskId)
	if Task.ID == 0 {
		return false
	}
	return true
}

func checkIfUserHasTaskInPeriod(task entities.Task) bool {
	PTask := []entities.Task{}

	result := database.Instance.Where("user_id = ?", task.UserID).Find(&PTask)
	if result.Error == nil && result.RowsAffected == 0 {
		return true
	}
	for i := 0; i < len(PTask); i++ {
		d1 := PTask[i].EndTime.Sub(task.StartTime).Seconds()
		d2 := PTask[i].StartTime.Sub(task.EndTime).Seconds()
		if d1*d2 < 0 {
			return false
		}
	}

	return true
}
