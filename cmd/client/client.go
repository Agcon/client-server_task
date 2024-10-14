package main

import (
	"client-server_task/db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"os"
	"sync"
)

var mu sync.Mutex
var serverAddress string = "http://localhost:8082"
var log = logrus.New()

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/getUsers", getUsersHandler).Methods("GET")
	router.HandleFunc("/getUsers/age", getUsersByAgeHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Info("Сервер запущен на http://localhost:8083")
	if err := http.ListenAndServe(":8083", router); err != nil {
		log.Fatal(err)
	}
}

func getUsersByAgeHandler(writer http.ResponseWriter, request *http.Request) {
	age := request.URL.Query().Get("value")
	log.WithFields(logrus.Fields{
		"age": age,
	}).Info("Получение пользователей по возрасту")

	users, err := getUsersFromServerByAge(age)
	if err != nil {
		log.WithError(err).Error("Ошибка получения данных с сервера по возрасту")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = saveRecordsToJSON(users)
	if err != nil {
		log.WithError(err).Error("Ошибка при сохранении данных в JSON")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Данные успешно записаны в JSON файл")
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}

func getUsersHandler(writer http.ResponseWriter, request *http.Request) {
	limit := request.URL.Query().Get("value")
	log.WithFields(logrus.Fields{
		"limit": limit,
	}).Info("Получение пользователей с ограничением")

	users, err := getUsersFromServer(limit)
	if err != nil {
		log.WithError(err).Error("Ошибка получения данных с сервера")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = saveRecordsToJSON(users)
	if err != nil {
		log.WithError(err).Error("Ошибка при сохранении данных в JSON")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Данные успешно записаны в JSON файл")
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}

func saveRecordsToJSON(users []db.User) error {
	mu.Lock()
	defer mu.Unlock()

	log.Info("Сохранение данных в JSON файл")
	file, err := os.OpenFile("users.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.WithError(err).Error("Ошибка при открытии файла")
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, user := range users {
		if err := encoder.Encode(user); err != nil {
			log.WithError(err).Error("Ошибка при записи пользователя в файл")
			return err
		}
		log.WithFields(logrus.Fields{
			"user": user,
		}).Info("Пользователь записан в файл")
	}
	log.Info("Все пользователи успешно записаны")
	return nil
}

func getUsersFromServerByAge(age string) ([]db.User, error) {
	url := fmt.Sprintf("%s/users/%s", serverAddress, age)
	log.WithFields(logrus.Fields{
		"url": url,
	}).Info("Запрос данных с сервера по возрасту")

	response, err := http.Get(url)
	if err != nil {
		log.WithError(err).Error("Ошибка при выполнении запроса к серверу по возрасту")
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.WithError(err).Error("Ошибка при чтении ответа сервера")
		return nil, err
	}

	var users []db.User
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.WithError(err).Error("Ошибка при декодировании данных с сервера")
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"count": len(users),
	}).Info("Успешно получены данные с сервера по возрасту")
	return users, nil
}

func getUsersFromServer(limit string) ([]db.User, error) {
	url := fmt.Sprintf("%s/users", serverAddress)
	log.WithFields(logrus.Fields{
		"url": url,
	}).Info("Запрос данных с сервера")

	response, err := http.Get(url)
	if err != nil {
		log.WithError(err).Error("Ошибка при выполнении запроса к серверу")
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.WithError(err).Error("Ошибка при чтении ответа сервера")
		return nil, err
	}

	var users []db.User
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.WithError(err).Error("Ошибка при декодировании данных с сервера")
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"count": len(users),
	}).Info("Успешно получены данные с сервера")
	return users, nil
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("cmd/client/templates/index.html")
	if err != nil {
		log.WithError(err).Error("Ошибка при загрузке шаблона")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Главная страница отображена")
	t.Execute(writer, nil)
}
