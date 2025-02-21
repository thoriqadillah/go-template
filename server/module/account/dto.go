package account

import "app/lib/validator"

type createUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6,lte=16"`
}

func (u *createUser) Validate() error {
	return validator.Struct(u)
}

type loginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6,lte=16"`
}

func (u *loginUser) Validate() error {
	return validator.Struct(u)
}
