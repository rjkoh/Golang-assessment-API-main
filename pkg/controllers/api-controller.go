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
	// struct to store data from json request body
	type Reg struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}
	var reg Reg
	err := utils.ParseBody(req, &reg)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	for _, student := range reg.Students {
		// add each student into the model and database
		err = models.AddStudent(reg.Teacher, student)
		if err != nil {
			handleError(writer, req, err)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func FindCommonStudents(writer http.ResponseWriter, req *http.Request) {
	// extract teachers from request URL
	teachers, ok := req.URL.Query()["teacher"]
	if !ok {
		handleError(writer, req, fmt.Errorf("Unable to read teacher parameters"))
		return
	}

	// find students common to the teachers specified
	rows, err := models.FindCommon(teachers)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	// struct to store the students specified to convert to json
	type Response struct {
		Students []string `json:"students"`
	}

	response := Response{Students: make([]string, 0)}
	// add all students found to the slice
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			handleError(writer, req, err)
			return
		}
		response.Students = append(response.Students, email)
	}

	// convert students slice into json
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
	// struct to store the student's email
	type Suspended struct {
		Student string `json:"student"`
	}
	var email Suspended
	err := utils.ParseBody(req, &email)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	err = models.Suspend(email.Student)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func RetrieveStudentsForNotification(writer http.ResponseWriter, req *http.Request) {
	// struct to store the data in the http request body
	type Noti struct {
		Teacher      string `json:"teacher"`
		Notification string `json:"notification"`
	}

	var rcvNoti Noti
	err := utils.ParseBody(req, &rcvNoti)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	// find students that are notifiable, either mentioned in the message or under the teacher
	rows, err := models.GetNotifiableStudents(getEmails(rcvNoti.Notification), rcvNoti.Teacher)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	// struct to store students found
	type Response struct {
		Recipients []string `json:"recipients"`
	}

	response := Response{Recipients: make([]string, 0)}
	// add all students into the slice to convert to json
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			handleError(writer, req, err)
			return
		}
		response.Recipients = append(response.Recipients, email)
	}

	// convert to json
	jsonData, err := json.Marshal(response)
	if err != nil {
		handleError(writer, req, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonData)
}

// extract email from the given notification
func getEmails(str string) []string {
	// with format starting with @, followed by alphanumeric characters and then a "." with alphanumeric characters
	re := regexp.MustCompile(`\b\w+@\w+\.\w+\b`)
	return re.FindAllString(str, -1)
}

// return a http bad request status with the error message
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
