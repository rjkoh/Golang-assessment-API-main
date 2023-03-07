package models_test

import (
	"testing"

	"github.com/rjkoh/golang-assessment-api/pkg/models"
)

func TestAddStudentSuccess(t *testing.T) {
	err := models.AddStudent("teacher@test.com", "student@test.com")
	if err != nil {
		t.Errorf("Error adding student with valid teacher and valid student" + err.Error())
	}
}

func TestAddStudentInvalidEmail(t *testing.T) {
	err := models.AddStudent("teacher@test.com", "")
	if err == nil {
		t.Errorf("Test Failed: allows invalid input for student")
	}

	err = models.AddStudent("teacher@test.com", "stud.com")
	if err == nil {
		t.Errorf("Test Failed: allows invalid input for student")
	}

	err = models.AddStudent("", "student@test.com")
	if err == nil {
		t.Errorf("Test Failed: allows invalid input for teacher")
	}

	err = models.AddStudent("teacher.com", "student@test.com")
	if err == nil {
		t.Errorf("Test Failed: allows invalid input for teacher")
	}
}

func TestFindCommonSuccess2Teachers(t *testing.T) {
	_, err := models.FindCommon([]string{"teacher1@gmail.com", "teacher2@gmail.com"})
	if err != nil {
		t.Errorf("Test Failed: Unable to find common students with 2 teachers" + err.Error())
	}
}

func TestFindCommonSuccess1Teacher(t *testing.T) {
	_, err := models.FindCommon([]string{"teacher1@gmail.com"})
	if err != nil {
		t.Errorf("Test Failed: Unable to find common students with 1 teacher" + err.Error())
	}
}

func TestFindCommonFailure(t *testing.T) {
	_, err := models.FindCommon([]string{})
	if err == nil {
		t.Errorf("Test Failed: Accepts 0 teachers")
	}
}

func TestSuspendSuccess(t *testing.T) {
	err := models.Suspend("student@test.com")
	if err != nil {
		t.Errorf("Test Failed: Unable to suspend valid student" + err.Error())
	}
}

func TestSuspendFailure(t *testing.T) {
	err := models.Suspend("student.com")
	if err == nil {
		t.Errorf("Test Failed: Accepts invalid student email")
	}

	err = models.Suspend("")
	if err == nil {
		t.Errorf("Test Failed: Accepts empty student email")
	}
}

func TestGetNotifiableStudentsSuccess(t *testing.T) {
	_, err := models.GetNotifiableStudents([]string{"student1@test.com", "student2@test.com"}, "teacher@test.com")
	if err != nil {
		t.Errorf("Test Failed: Unable to get notifiable students with valid students and teacher" + err.Error())
	}
}

func TestGetNotifiableStudentsSuccessNoEmails(t *testing.T) {
	_, err := models.GetNotifiableStudents([]string{}, "teacher@test.com")
	if err != nil {
		t.Errorf("Test Failed: Unable to get notifiable students with valid students and teacher" + err.Error())
	}
}

func TestGetNotifiableStudentsInvalidTeacher(t *testing.T) {
	_, err := models.GetNotifiableStudents([]string{"student1@test.com", "student2@test.com"}, "teacher")
	if err == nil {
		t.Errorf("Test Failed: GetNotifiable accepts Invalid Teacher email")
	}
}

func TestGetNotifiableStudentsInvalidEmail(t *testing.T) {
	_, err := models.GetNotifiableStudents([]string{"student1@test.com", "student2"}, "teacher@test.com")
	if err == nil {
		t.Errorf("Test Failed: GetNotifiable accepts Invalid Student mentioned (@-ed) email")
	}
}
