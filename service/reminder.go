package service

import (
	"PersonalTaskManagement/database"
	"PersonalTaskManagement/entities"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)

var From string
var Password string
var RemindList = make(map[int]entities.Task)

func ReminderScheduller() {
	for {
		log.Println("ReminderScheduller ", RemindList)
		var Tasks []entities.Task
		database.Instance.Find(&Tasks)
		for _, task := range Tasks {
			diffStart := time.Now().Sub(task.StartTime).Seconds()
			if diffStart > 0 {
				diff := time.Now().Sub(task.UpdatedAt).Hours()
				if diff > float64(task.Reminder) {
					RemindList[int(task.ID)] = task
					if len(From) > 0 {
						RemindByEmail(task)
					}
				}
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
	var User entities.User
	database.Instance.First(&User, task.UserID)

	EmailShouldSend := false
	EmailSent := entities.EmailsReminder{}
	database.Instance.Where("task_id = ?", task.ID).Find(&EmailSent)
	if EmailSent.ID == 0 {
		EmailSent.TaskId = int(task.ID)
		EmailSent.Email = User.Email
		database.Instance.Create(&EmailSent)
		EmailShouldSend = true
	}

	diff := time.Now().Sub(task.UpdatedAt).Hours()
	if diff > float64(task.Reminder) {
		EmailShouldSend = true
	}

	if EmailShouldSend {
		SendEmail(User.Email, task)
		database.Instance.Save(&EmailSent)
	}
}

func SendEmail(e string, t entities.Task) {

	// Receiver email address.
	to := []string{e}
	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("task reminder:" + t.Description)

	// Authentication.
	auth := smtp.PlainAuth("", From, Password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, From, to, message)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Email Sent Successfully!")
}
