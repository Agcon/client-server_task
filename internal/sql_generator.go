package main

import (
	"client-server_task/db"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	database := &db.DB{}
	err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatalf("Could not close database connection: %v", err)
		}
	}()

	cities := [4]string{"Moscow", "Smolensk", "Kazan", "Kaluga"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10000; i++ {
		user := db.User{
			Name:         fmt.Sprintf("Name %d", i+1),
			LastName:     fmt.Sprintf("Lastname %d", (i+1)*10),
			Age:          r.Intn(60) + 10,
			PlaceOfBirth: cities[r.Intn(4)],
			Salary:       r.Intn(1000000) + 40000,
		}
		if err := database.InsertRecord(&user); err != nil {
			log.Fatalf("Could not create a new record: %v", err)
		}
	}
}
