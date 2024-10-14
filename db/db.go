package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Conn *gorm.DB
}

type User struct {
	ID           uint   `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"lastName"`
	Age          int    `json:"age"`
	PlaceOfBirth string `json:"placeOfBirth"`
	Salary       int    `json:"salary"`
}

type Database interface {
	Connect() error
	Close() error
	FetchAllRecords() ([]User, error)
	FetchRecords(limit int) ([]User, error)
	FetchRecordsByAge(age int) ([]User, error)
	InsertRecord(user *User) error
}

func (db *DB) Connect() error {
	dsn := "host=localhost user=postgres password=admin dbname=client-server_db port=5432 sslmode=disable"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	if err := conn.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("could not migrate database: %w", err)
	}
	db.Conn = conn
	return nil
}

func (db *DB) Close() error {
	sqlDB, err := db.Conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *DB) FetchAllRecords() ([]User, error) {
	var users []User
	result := db.Conn.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (db *DB) FetchRecords(limit int) ([]User, error) {
	var users []User
	result := db.Conn.Limit(limit).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (db *DB) FetchRecordsByAge(age int) ([]User, error) {
	var users []User
	result := db.Conn.Where("age >= ?", age).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (db *DB) InsertRecord(user *User) error {
	result := db.Conn.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
