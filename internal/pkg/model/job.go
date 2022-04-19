package model

type Job struct {
	Name    string `json:"name" binding:"required,min=3,max=255"`
	Command string `json:"command" binding:"required"`
	Expr    string `json:"expr" binding:"required"`
}
