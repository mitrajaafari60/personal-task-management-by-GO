package entities

type User struct {
	ID        int
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"unique;not null;type:varchar(100);default:null"`
}
