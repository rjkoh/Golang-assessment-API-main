package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/rjkoh/golang-assessment-api/pkg/config"
)

type Student struct {
	Email       string `form:"email" json:"email"`
	Teacher     string `form:"teacher" json:"teacher"`
	IsSuspended bool   `form:"isSuspended" json:"isSuspended"`
}

func AddStudent(teacher string, email string) error {
	config.Connect()
	db = config.GetDB()
	defer db.Close()
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(teacher) || !re.MatchString(email) {
		return fmt.Errorf("Invalid teacher email or student email")
	}

	query := "INSERT INTO Student(email, teacher, isSuspended) VALUES (?, ?, ?)"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, email, teacher, false)
	if err != nil {
		log.Printf("Error %s when executing", err)
		return err
	}

	return nil
}

func FindCommon(teachers []string) (*sql.Rows, error) {
	config.Connect()
	db = config.GetDB()
	defer db.Close()
	if len(teachers) == 0 {
		return nil, fmt.Errorf("Missing teachers parameters")
	}
	query := "SELECT DISTINCT email FROM Student WHERE teacher = '" + teachers[0] + "'"

	for i := 1; i < len(teachers); i++ {
		intersect := "INTERSECT SELECT DISTINCT email FROM Student WHERE teacher = '" + teachers[i] + "'"
		query = query + " " + intersect
	}

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func Suspend(email string) error {
	config.Connect()
	db = config.GetDB()
	defer db.Close()
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(email) {
		return fmt.Errorf("Invalid teacher email or student email")
	}

	query := "UPDATE Student SET isSuspended = TRUE WHERE email = ?"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, email)
	if err != nil {
		log.Printf("Error %s when executing", err)
		return err
	}

	return nil
}

func GetNotifiableStudents(emails []string, teacher string) (*sql.Rows, error) {
	config.Connect()
	db = config.GetDB()
	defer db.Close()
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(teacher) {
		return nil, fmt.Errorf("Invalid teacher email")
	}

	withTeacherQuery := fmt.Sprintf("SELECT DISTINCT email FROM Student WHERE teacher = '%s' AND isSuspended = FALSE", teacher)

	var mentionedQuery string
	var query string
	if len(emails) > 0 {
		mentionedQuery = "SELECT DISTINCT email FROM Student WHERE email = '" + emails[0] + "'"

		for i := 1; i < len(emails); i++ {
			if !re.MatchString(emails[i]) {
				return nil, fmt.Errorf("Invalid student email")
			}
			next := "OR email = '" + emails[i] + "'"
			mentionedQuery += next
		}
		query = withTeacherQuery + " UNION " + mentionedQuery
	} else {
		mentionedQuery = ""
		query = withTeacherQuery
	}

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	return rows, nil
}
