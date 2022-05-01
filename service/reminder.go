package service

import (
	"PersonalTaskManagement/database"
	"PersonalTaskManagement/entities"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

var RemindList = make(map[int]entities.Task)

func ReminderScheduller() {
	for {
		log.Println("ReminderScheduller ", RemindList)
		var Tasks []entities.Task
		database.Instance.Find(&Tasks)
		for _, task := range Tasks {
			diff := time.Now().Sub(task.UpdatedAt).Seconds()
			if diff > float64(task.Reminder) {
				RemindList[int(task.ID)] = task
				RemindByEmail(task)
			}
		}
		time.Sleep(20 * time.Second)
		for _, remind := range RemindList {
			delete(RemindList, int(remind.ID))
		}
	}
}

func SiteRemindTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RemindList)
}

func SiteRemindTaskVisited(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var Task entities.Task
	json.NewDecoder(r.Body).Decode(&Task)
	err := UpdateRemind(strconv.Itoa(int(Task.ID)))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode("Task Reminder Updated!")
}

func UpdateRemind(TaskId string) error {
	var Task entities.Task
	result := database.Instance.First(&Task, TaskId)
	if result.RowsAffected == 0 {
		return errors.New("Task Not Found!")
	}
	Task.RemindCounts = Task.RemindCounts + 1
	result = database.Instance.Save(&Task)
	if result.Error != nil {
		return errors.New(result.Error.Error())
	}
	return nil
}

func RemindByEmail(task entities.Task) {

}