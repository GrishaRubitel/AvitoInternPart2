package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID        string    `gorm:"type:varchar(100);primaryKey"`
	Username  string    `gorm:"type:varchar(50);unique;not null"`
	FirstName string    `gorm:"type:varchar(50)"`
	LastName  string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (Employee) TableName() string {
	return "employee"
}

func EmployeeMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Employee{})
	if err != nil {
		return err
	}
	return nil
}

func FindUserInTable(db *gorm.DB, username string) (Employee, error) {
	var emp Employee
	query := db.Where("username = ?", username)
	result := query.Find(&emp)
	if result.Error != nil {
		return Employee{}, errors.New("error while searching for user")
	} else {
		return emp, nil
	}
}
