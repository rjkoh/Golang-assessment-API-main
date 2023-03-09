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

// to store the student data
type Student struct {
	Email       string `form:"email" json:"email"`
	Teacher     string `form:"teacher" json:"teacher"`
	IsSuspended bool   `form:"isSuspended" json:"isSuspended"`
}

func AddStudent(teacher string, email string) error {
	// connect to the database and close at the end of function
	config.Connect()
	db = config.GetDB()
	defer db.Close()

	// check that the given student and teacher emails are of valid format
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(teacher) || !re.MatchString(email) {
		return fmt.Errorf("Invalid teacher email or student email")
	}

	// create sql query
	query := "INSERT INTO Student(email, teacher, isSuspended) VALUES (?, ?, ?)"

	// create context with timeout
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelfunc()

	// prepare the database for query (prevent sql injection)
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing", err)
		return err
	}

	// execute query
	_, err = stmt.ExecContext(ctx, email, teacher, false)
	if err != nil {
		log.Printf("Error %s when executing", err)
		return err
	}

	return nil
}

func FindCommon(teachers []string) (*sql.Rows, error) {
	// connect to the database and close at the end of function
	config.Connect()
	db = config.GetDB()
	defer db.Close()

	// check that teacher(s) are provided
	if len(teachers) == 0 {
		return nil, fmt.Errorf("Missing teachers parameters")
	}

	// check that teacher email is of valid format
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(teachers[0]) {
		return nil, fmt.Errorf("Invalid teacher email")
	}

	// create query
	query := "SELECT DISTINCT email FROM Student WHERE teacher = '" + teachers[0] + "'"

	// append any remaining valid teacher emails
	for i := 1; i < len(teachers); i++ {
		if !re.MatchString(teachers[i]) {
			return nil, fmt.Errorf("Invalid teacher email")
		}
		intersect := "INTERSECT SELECT DISTINCT email FROM Student WHERE teacher = '" + teachers[i] + "'"
		query = query + " " + intersect
	}

	// query the database
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func Suspend(email string) error {
	// connect to the database and close at the end of function
	config.Connect()
	db = config.GetDB()
	defer db.Close()

	// check that given student email is of valid format
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(email) {
		return fmt.Errorf("Invalid student email")
	}

	// create query
	query := "UPDATE Student SET isSuspended = TRUE WHERE email = ?"

	// create context with timeout
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelfunc()

	// prepare the db for the query to prevent sql injection
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing", err)
		return err
	}

	// execute query with given data
	_, err = stmt.ExecContext(ctx, email)
	if err != nil {
		log.Printf("Error %s when executing", err)
		return err
	}

	return nil
}

func GetNotifiableStudents(emails []string, teacher string) (*sql.Rows, error) {
	// connect to the database and close at the end of function
	config.Connect()
	db = config.GetDB()
	defer db.Close()

	// check teacher email is of valid format
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(teacher) {
		return nil, fmt.Errorf("Invalid teacher email")
	}

	// create query for students under the teacher that are not suspended
	withTeacherQuery := fmt.Sprintf("SELECT DISTINCT email FROM Student WHERE teacher = '%s' AND isSuspended = FALSE", teacher)

	var mentionedQuery string
	var query string
	// create query for students that are mentioned in the notification
	if len(emails) > 0 {
		mentionedQuery = "SELECT DISTINCT email FROM Student WHERE email = '" + emails[0] + "'"

		// check student emails are of valid format and append them to the query
		for i := 1; i < len(emails); i++ {
			if !re.MatchString(emails[i]) {
				return nil, fmt.Errorf("Invalid student email")
			}
			next := "OR email = '" + emails[i] + "'"
			mentionedQuery += next
		}

		// union all valid rows for the conditions
		query = withTeacherQuery + " UNION " + mentionedQuery
	} else {
		mentionedQuery = ""
		query = withTeacherQuery
	}

	// execute query
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	return rows, nil
}
