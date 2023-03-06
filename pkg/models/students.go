package models

type Student struct {
	email   string `form:"email" json:"email"`
	teacher string `form:"teacher" json:"teacher"`
}
