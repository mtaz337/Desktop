package dto

import validation "github.com/go-ozzo/ozzo-validation"

type SendSms struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

func (s SendSms) Validate() error {
	return validation.ValidateStruct(&s, 
		validation.Field(&s.Number, validation.Required),
		validation.Field(&s.Message, validation.Required, validation.Length(5, 300)),
	)
}