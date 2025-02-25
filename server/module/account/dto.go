package account

import (
	"app/lib/validator"
	"time"
)

type createUser struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required,gte=6,lte=16"`
	Source   string  `json:"source" validate:"required,oneof='email google"`
}

func (u *createUser) Validate() error {
	return validator.Struct(u)
}

type login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6,lte=16"`
}

func (u *login) Validate() error {
	return validator.Struct(u)
}

type updateUser struct {
	VerifiedAt time.Time
	Password   string `json:"password" validate:"gte=6,lte=16"`
	Name       string `json:"name"`
}

type changePassword struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6,lte=16"`
	Otp      string `json:"otp" validate:"required"`
}

func (u *changePassword) Validate() error {
	return validator.Struct(u)
}
