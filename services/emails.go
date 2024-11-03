package services

import (
	"bytes"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/repository"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/hegdeshashank73/glamr-backend/vendors"
	"github.com/sirupsen/logrus"
)

var CharSet = "UTF-8"

func GetEmailTemplate(req entities.GetEmailTemplateReq) (entities.GetEmailTemplateRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.GetEmailTemplate", st)

	res := entities.GetEmailTemplateRes{}

	ret, derr := repository.GetEmailTemplate(vendors.DBMono, entities.GetEmailTemplateArg(req))
	if derr != nil {
		return res, derr
	}
	res.Template = ret.Template
	res.Message.Title = "Email template created successfully"
	return res, nil
}

func SendMagiclink(recipientEmail string, data entities.DataSendMagiclink) errors.GlamrError {
	ret, derr := GetEmailTemplate(entities.GetEmailTemplateReq{Name: "AuthMagiclink"})
	if derr != nil {
		return derr
	}
	return sendEmailToEmail2(
		recipientEmail,
		ret.Template.Subject,
		ret.Template.Body,
		data)
}

func sendEmailToEmail2(email string, subjectTemplate string, bodyTemplate string, data any) errors.GlamrError {
	senderEmail := "Glamr <no-reply@glamr.us>"
	replyToEmail := "shashin.bhaskar@gmail.com"
	tmpl, err := template.New("bodyTemplate").Parse(bodyTemplate)
	if err != nil {
		logrus.Error("failed to parse template,", err)
		return errors.GlamrErrorGeneralServerError("Oops! there's a glitch in the matric and we could not send you magiclink at the moment.")
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		logrus.Error("failed to render template,", err)
		return errors.GlamrErrorGeneralServerError("Oops! there's a glitch in the matric and we could not send you magiclink at the moment.")
	}

	tmplSubject, err := template.New("subjectTemplate").Parse(subjectTemplate)
	if err != nil {
		logrus.Error("failed to parse template,", err)
		return errors.GlamrErrorGeneralServerError("Oops! there's a glitch in the matric and we could not send you magiclink at the moment.")
	}

	var subject bytes.Buffer
	err = tmplSubject.Execute(&subject, data)
	if err != nil {
		logrus.Error("failed to render template,", err)
		return errors.GlamrErrorGeneralServerError("Oops! there's a glitch in the matric and we could not send you magiclink at the moment.")
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject.String()),
			},
		},
		Source: aws.String(senderEmail),
		ReplyToAddresses: []*string{
			aws.String(replyToEmail),
		},
	}

	_, err = vendors.SESClient.SendEmail(input)
	if err != nil {
		logrus.Error(err)
		return errors.GlamrErrorGeneralServerError("Oops! there's a glitch in the matric and we could not send you magiclink at the moment.")
	}
	return nil
}

func CreateEmailTemplate(req entities.CreateEmailTemplateReq) (entities.CreateEmailTemplateRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.CreateEmailTemplate", st)

	res := entities.CreateEmailTemplateRes{}

	derr := repository.CreateEmailTemplate(vendors.DBMono, entities.CreateEmailTemplateArg(req))
	if derr != nil {
		return res, derr
	}
	res.Message.Title = "Email template created successfully"
	return res, nil
}
