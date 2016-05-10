package model

// Category table
type Category struct {
	Id int64

	ConfigId int64

	Name        string
	Description string

	Color string
}
