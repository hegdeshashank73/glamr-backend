package repository

import (
	"time"

	"github.com/hegdeshashank73/glamr-backend/dsql"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/sirupsen/logrus"
)

func GetEmailTemplate(tx dsql.Tx, arg entities.GetEmailTemplateArg) (entities.GetEmailTemplateRet, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("repository.GetEmailTemplate", st)

	ret := entities.GetEmailTemplateRet{}
	err := tx.QueryRow(
		`select id, name, body, subject from templates_emails where name = ?`, arg.Name).Scan(
		&ret.Template.ID, &ret.Template.Name, &ret.Template.Body, &ret.Template.Subject,
	)

	if err != nil {
		logrus.Error(err)
		return ret, errors.GlamrErrorDatabaseIssue()
	}

	return ret, nil
}
