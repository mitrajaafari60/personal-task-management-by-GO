# this is a sample project based on Go for Task Mangment APIs


## Introduction

This Application perform the requests as a personal task management app bata:
Minimal requirements:
• Create user
  ◦ Minimal set of properties:
      ▪ email----------> Email validation and being unique in users table is checked
      ▪ first name
      ▪ last name
• Delete user
• Get user
• Get users
• Create a task
    ◦ task has to be assigned to a user   ---------> user Validation is Checked
    ◦ task can occupy some certain amount of time-->  user Validation and start/End time of the task checked
    ◦ task can’t overlap with other tasks ----------> start/end time of task of the user checked with his/her other tasks in the Database not to be overlapped
    ◦ reminder period-------------------------------> default is set to 24Hour , but it can be defined in task creation
• Delete task
• Get task
• Get tasks
  ◦ get list of the task
• Reminder notification
  ◦ by email  -------> in service/reminder.go a goroutin checks the task and their timing if notification is needed it fills in a slice of map which is Get by web Site,and Emails and messages by TaskId is saved on Database to check timing and perevent multiple Email sending.
  ◦ on-site  --------> as described reminderList is saved in list by goroutin and each time the method(GET /api/reminders) called it respond the array of tasks which should be done , and also (POST /api/reminders) method is developed, this method updates the task time UpdatedAt for the spesific taskId ,So the notification of the task on site can ignored  and not to be annoy the user each time  untill next remind period.
* ### *Description:*
    * This project  backend-server is based on Rest webservice with Postgres Database
## Getting Started
Clone the git repository in your system and then cd into the project root directory

```bash
$ git clone https://github.com/mitrajaafari60/personal-task-management-by-GO 
$ cd personal-task-management-by-GO 
```
### Project Structure

* Project skeleton is as below:
    * root of project   :
        * the main part of the project for setting up and HTTP server is launched in this package.
        * config.json is a config file with 4 needed data first is connection_string for DB connection, you can change with your own configuration, listen to Port and The Email, and password for sending notifications as mentioned in assignment criteria. (in case Email was Null nothing will be sent by Email but list of notofications can captured with GET /api/reminders)
        * controller
            * <sub>functions in  this folder controls the request formate and in cas of success, service is called and the appropriate response is generated and will be sent back to the client </sub>
        * entities
            * <sub>In this folder,data model for restApi and DataBase included.</sub> 
        * service
            * <sub>all related parts of the service such as reminder with getting task From database and checking their start and stop time and also reminder duration is done in service package.</sub>
        * database
            * <sub>general functions which developed and used in the project is categorized in pkg folder</sub>

### how to lunch
  for easy use of the project for those who don't have a Go compiler in their system 
  they can just run simply as below, also lunch http server on port 8080 will be lunched.
  for running by go it can be run by main which is located in root folder with go run main.go ...
   
### how to check        
  it is so simple!you can run and test the method with the postman collection included in the base folder named as `taskManagement.postman_collection.postman_collection.json`  
  