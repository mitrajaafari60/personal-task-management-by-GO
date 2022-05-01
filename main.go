package main

import (
	"PersonalTaskManagement/controllers"
	"PersonalTaskManagement/database"
	"PersonalTaskManagement/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/spf13/viper"
)

type Config struct {
	Port             string `mapstructure:"port"`
	ConnectionString string `mapstructure:"connection_string"`
	Email            string `mapstructure:"email"`
	Password         string `mapstructure:"password"`
}

var AppConfig *Config
var DB *gorm.DB

func main() {

	// Load Configurations from config.json using Viper
	LoadAppConfig()

	// Initialize Database
	database.Connect(AppConfig.ConnectionString)
	database.Migrate()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	// run Reminder goroutin to Remind Users
	service.From = AppConfig.Email
	service.Password = AppConfig.Password
	go service.ReminderScheduller()
	// Register Routes
	RegisterRoutes(router)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))
}

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/users", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", controllers.GetUserById).Methods("GET")
	router.HandleFunc("/api/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", controllers.DeleteUser).Methods("DELETE")

	router.HandleFunc("/api/tasks", controllers.GetTasks).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", controllers.GetTaskById).Methods("GET")
	router.HandleFunc("/api/tasks", controllers.CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks/{id}", controllers.DeleteTask).Methods("DELETE")

	router.HandleFunc("/api/reminders", service.SiteRemindTasks).Methods("GET")
	router.HandleFunc("/api/reminders", service.SiteRemindTaskVisited).Methods("POST")

}

func LoadAppConfig() {
	log.Println("Loading Server Configurations...")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal(err)
	}
}
