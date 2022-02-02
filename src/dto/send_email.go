package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SendEmail struct {
	From         string      `json:"from"`
	To           string      `json:"to"`
	Subject      string      `json:"subject"`
	Message      string      `json:"message"`
	TemplateCode string      `json:"template_code"`
	Data         interface{} `json:"data"`
	CC           *string     `json:"cc"`
	BCC          *string     `json:"bcc"`
}

func (s SendEmail) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.From, validation.Required, is.Email),
		validation.Field(&s.To, validation.Required, is.Email),
		validation.Field(&s.CC, is.Email),
		validation.Field(&s.BCC, is.Email),
		validation.Field(&s.Subject, validation.Required),
		validation.Field(&s.Message, validation.Required),
		validation.Field(&s.TemplateCode, validation.Required),
		validation.Field(&s.Data, validation.Required),
	)
}
