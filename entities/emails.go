package entities

import "github.com/hegdeshashank73/glamr-backend/common"

type DataSendMagiclink struct {
	Magiclink string
}

type GetEmailTemplateReq struct {
	Name string `json:"name"`
}

type GetEmailTemplateRes struct {
	Template EmailTemplate `json:"template"`
	Message  Message       `json:"message"`
}

type GetEmailTemplateRet GetEmailTemplateRes

type EmailTemplate struct {
	Subject string           `json:"subject"`
	Body    string           `json:"body"`
	Name    string           `json:"name"`
	ID      common.Snowflake `json:"id"`
}

type GetEmailTemplateArg GetEmailTemplateReq
