package example

import (
	"app/lib/validator"
)

type Foo struct {
	Name string `json:"name" validate:"required"`
}

func (f *Foo) Validate() error {
	return validator.Struct(f)
}
