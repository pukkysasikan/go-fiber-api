package models

type Category struct {
	Name string `json:"name,omitempty" validate:"required"`
}
