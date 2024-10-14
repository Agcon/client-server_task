package main

import (
	"client-server_task/db"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var database = &db.DB{}

func main() {
	err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatalf("Could not close database connection: %v", err)
		}
	}()
	router := mux.NewRouter()

	router.HandleFunc("/users", getUsersHandler).Methods("GET")
	router.HandleFunc("/users/{age}", getUsersByAge).Methods("GET")

	log.Println("Server is running...")

	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func getUsersByAge(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	ageStr := vars["age"]
	age, err := strconv.Atoi(ageStr)
	if err != nil || age <= 10 {
		http.Error(writer, "Invalid age parameter", http.StatusBadRequest)
		return
	}

	users, err := database.FetchRecordsByAge(age)

	if err != nil {
		http.Error(writer, "Error fetching records by age", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)

	if err != nil {
		http.Error(writer, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

func getUsersHandler(writer http.ResponseWriter, request *http.Request) {
	limitStr := request.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 0
	}

	var users []db.User
	if limit == 0 {
		users, err = database.FetchAllRecords()
	} else {
		users, err = database.FetchRecords(limit)
	}

	if err != nil {
		http.Error(writer, "Error fetching records", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)

	if err != nil {
		http.Error(writer, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}
