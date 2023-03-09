package controllers

import (
	"encoding/json"
	"fmt"
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
			handleError(writer, req, err)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)

}

func FindCommonStudents(writer http.ResponseWriter, req *http.Request) {
	teachers, ok := req.URL.Query()["teacher"]
	if !ok {
		handleError(writer, req, fmt.Errorf("Unable to read teacher parameters"))
		return
	}

	rows, err := models.FindCommon(teachers)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	type Response struct {
		Students []string `json:"students"`
	}

	response := Response{Students: make([]string, 0)}
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			handleError(writer, req, err)
			return
		}
		response.Students = append(response.Students, email)
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
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
		handleError(writer, req, err)
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
	rows, err := models.GetNotifiableStudents(getEmails(rcvNoti.Notification), rcvNoti.Teacher)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	type Response struct {
		Recipients []string `json:"recipients"`
	}

	response := Response{Recipients: make([]string, 0)}
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			handleError(writer, req, err)
			return
		}
		response.Recipients = append(response.Recipients, email)
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		handleError(writer, req, err)
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

func handleError(writer http.ResponseWriter, req *http.Request, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	jsondata, er := json.Marshal(struct {
		Message string `json:"message"`
	}{Message: err.Error()})
	if er != nil {
		log.Printf("failed to write response: %v", er)
	}
	writer.Write(jsondata)
}
