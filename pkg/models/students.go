package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Student struct {
	Email       string `form:"email" json:"email"`
	Teacher     string `form:"teacher" json:"teacher"`
	IsSuspended bool   `form:"isSuspended" json:"isSuspended"`
}

func AddStudent(teacher string, email string) error {
	query := "INSERT INTO Student(email, teacher, isSuspended) VALUES (?, ?, ?)"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing", err)
		return err
	}

	res, err := stmt.ExecContext(ctx, email, teacher, false)
	if err != nil {
		log.Printf("Error %s when executing", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding students affected", err)
		return err
	}

	log.Printf("%d student created, %+v", rows, email)

	return nil
}

func FindCommon(teachers []string) *sql.Rows {
	query := "SELECT email FROM Student WHERE teacher = '" + teachers[0] + "'"

	for i := 1; i < len(teachers); i++ {
		intersect := "INTERSECT SELECT email FROM Student WHERE teacher = '" + teachers[i] + "'"
		query = query + " " + intersect
	}

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return rows
}

func Suspend(email string) error {
	query := "UPDATE Student SET isSuspended = TRUE WHERE email = ?"

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing", err)
		return err
	}

	res, err := stmt.ExecContext(ctx, email)
	if err != nil {
		log.Printf("Error %s when executing", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding students affected", err)
		return err
	}

	log.Printf("%d student suspended, %s", rows, email)

	return nil
}

func GetNotifiableStudents(emails []string, teacher string) *sql.Rows {
	withTeacherQuery := fmt.Sprintf("SELECT email FROM Student WHERE teacher = '%s' AND isSuspended = FALSE", teacher)

	mentionedQuery := "SELECT email FROM Student WHERE email = '" + emails[0] + "'"

	for i := 1; i < len(emails); i++ {
		next := "OR email = '" + emails[i] + "'"
		mentionedQuery += next
	}

	query := withTeacherQuery + " UNION " + mentionedQuery

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return rows
}
