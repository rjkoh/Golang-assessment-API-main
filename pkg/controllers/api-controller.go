package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/rjkoh/golang-assessment-api/pkg/models"
	"github.com/rjkoh/golang-assessment-api/pkg/utils"
)

func RegisterStudents(writer http.ResponseWriter, req *http.Request) {
	type Reg struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}
	var reg Reg
	utils.ParseBody(req, &reg)
	for _, student := range reg.Students {
		err := models.AddStudent(reg.Teacher, student)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			er := json.NewEncoder(writer).Encode(struct{ message string }{message: err.Error()})
			if er != nil {
				log.Printf("failed to write response: %v", er)
			}
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)

}

func FindCommonStudents(writer http.ResponseWriter, req *http.Request) {
	teachers, ok := req.URL.Query()["teacher"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		er := json.NewEncoder(writer).Encode(struct{ message string }{message: "Missing teacher parameter"})
		if er != nil {
			log.Printf("failed to write response: %v", er)
		}
		return
	}

	rows := models.FindCommon(teachers)

	type Response struct {
		Students []string `json:"students"`
	}

	response := Response{Students: make([]string, 0)}
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			er := json.NewEncoder(writer).Encode(struct{ message string }{message: err.Error()})
			if er != nil {
				log.Printf("failed to write response: %v", er)
			}
			return
		}
		response.Students = append(response.Students, email)
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		er := json.NewEncoder(writer).Encode(struct{ message string }{message: err.Error()})
		if er != nil {
			log.Printf("failed to write response: %v", er)
		}
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}

func SuspendStudent(writer http.ResponseWriter, req *http.Request) {
	type Suspended struct {
		Student string `json:"student"`
	}
	var email Suspended
	utils.ParseBody(req, &email)
	err := models.Suspend(email.Student)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		er := json.NewEncoder(writer).Encode(struct{ message string }{message: err.Error()})
		if er != nil {
			log.Printf("failed to write response: %v", er)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func RetrieveStudentsForNotification(writer http.ResponseWriter, req *http.Request) {
	type Noti struct {
		Teacher      string `json:"teacher"`
		Notification string `json:"notification"`
	}

	var rcvNoti Noti
	utils.ParseBody(req, &rcvNoti)
	rows := models.GetNotifiableStudents(getEmails(rcvNoti.Notification), rcvNoti.Teacher)

	type Response struct {
		Recipients []string `json:"recipients"`
	}

	response := Response{Recipients: make([]string, 0)}
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			er := json.NewEncoder(writer).Encode(struct{ message string }{message: err.Error()})
			if er != nil {
				log.Printf("failed to write response: %v", er)
			}
			return
		}
		response.Recipients = append(response.Recipients, email)
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		er := json.NewEncoder(writer).Encode(struct{ message string }{message: err.Error()})
		if er != nil {
			log.Printf("failed to write response: %v", er)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}

func getEmails(str string) []string {
	re := regexp.MustCompile(`\b\w+@\w+\.\w+\b`)
	return re.FindAllString(str, -1)
}
