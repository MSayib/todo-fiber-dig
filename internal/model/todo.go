package model

import "gorm.io/gorm"

// Todo merepresentasikan data di database
type Todo struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	IsDone      bool   `json:"is_done" gorm:"default:false"`
}

// ... request structs tetap sama
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required,min=10"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	IsDone      *bool  `json:"is_done,omitempty"`
}