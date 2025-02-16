package repository

import (
	"time"

	"github.com/hegdeshashank73/glamr-backend/common"
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
		`Select id, name, body, subject from templates_emails where name = $1`, arg.Name).Scan(
		&ret.Template.ID, &ret.Template.Name, &ret.Template.Body, &ret.Template.Subject,
	)

	if err != nil {
		logrus.Error(err)
		return ret, errors.GlamrErrorDatabaseIssue()
	}

	return ret, nil
}

func CreateEmailTemplate(tx dsql.Tx, arg entities.CreateEmailTemplateArg) errors.GlamrError {
	st := time.Now()
	defer utils.LogTimeTaken("repository.CreateEmailTemplate", st)

	id := common.GenerateSnowflake()

	query := `INSERT INTO templates_emails (id, name, body, subject)
          VALUES ($1, $2, $3, $4)
          ON CONFLICT (name) DO UPDATE SET
              name = EXCLUDED.name,
              body = EXCLUDED.body,
              subject = EXCLUDED.subject;`
	_, err := tx.Exec(query, id, arg.Name, arg.Body, arg.Subject)

	if err != nil {
		logrus.Error(err)
		return errors.GlamrErrorDatabaseIssue()
	}

	return nil
}
