package service

import (
	"github.com/go-playground/validator/v10"
)

type ValidateService struct {
	validate *validator.Validate
}

func NewValidateService() *ValidateService {
	return &ValidateService{validate: validator.New()}
}

func (s *ValidateService) validateStruct(data interface{}) error {
	return s.validate.Struct(data)
}
