package models

type User struct {
	ID    int
	Name  string
	Email string
}

type CreateUserInput struct {
	Name  string
	Email string
}
