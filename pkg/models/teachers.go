package models

type Teacher struct {
	email   string `form:"email" json:"email"`
	student string `form:"student" json:"student"`
}
