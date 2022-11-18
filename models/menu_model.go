package models

type Menu struct {
	Name        string `json:"name,omitempty" validate:"required"`
	Category    string `json:"category,omitempty" validate:"required"`
	Price       string `json:"price,omitempty" validate:"required"`
	Description string `json:"description,omitempty" validate:"required"`
	Image       string `json:"image,omitempty" validate:"required"`
}
