package schemas

import "github.com/lib/pq"

type UserCreateSchema struct {
	Name            string         `json:"name" validate:"required,max=25"`
	Username        string         `json:"username" validate:"required,max=16"` //alphanum+_
	Email           string         `json:"email" validate:"required,email"`
	Password        string         `json:"password" validate:"required,min=8"`
	ConfirmPassword string         `json:"confirmPassword" validate:"required,min=8"`
	Bio             string         `json:"bio" validate:"max=500"`
	Links           pq.StringArray `json:"links"`
}

type UserUpdateSchema struct {
	Name  *string         `json:"name" validate:"max=25"`
	Bio   *string         `json:"bio" validate:"max=500"`
	Links *pq.StringArray `json:"links"`
}

type UserLoginSchema struct {
	Username string `json:"username"`
	Password string `json:"password"`
}