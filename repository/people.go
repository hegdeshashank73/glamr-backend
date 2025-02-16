package repository

import (
	"database/sql"
	"time"

	"github.com/hegdeshashank73/glamr-backend/dsql"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/sirupsen/logrus"
)

func GetPerson(tx dsql.Tx, person *entities.Person, arg entities.GetPersonArg) (entities.GetPersonRet, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("repository.GetPerson", st)

	ret := entities.GetPersonRet{}
	var query string
	var args []any
	if arg.Token != "" {
		query = `SELECT id, first_name, last_name 
		FROM people_people pp
		INNER JOIN auth_tokens at ON at.user_id = pp.id where at.token = $1;`
		args = append(args, arg.Token)
	} else if arg.UserID > 0 {
		query = `SELECT id, first_name, last_name FROM people_people where id = $1;`
		args = append(args, arg.UserID)
	}

	err := tx.QueryRow(query, args...).Scan(
		&ret.Person.Id,
		&ret.Person.About.FirstName,
		&ret.Person.About.LastName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ret, errors.GlamrErrorGeneralNotFound("User not found")
		}
		logrus.Error(err)
		return ret, errors.GlamrErrorInternalServerError()
	}

	return ret, nil
}
func CreatePerson(tx *sql.Tx, arg *entities.CreatePersonArg) errors.GlamrError {
	_, err := tx.Exec(
		`INSERT INTO people_people (id, first_name, last_name) 
          VALUES ($1, $2, $3)
          ON CONFLICT (id) DO NOTHING;`,
		arg.Id, arg.FirstName, arg.LastName,
	)
	if err != nil {
		logrus.Error(err)
		return errors.GlamrErrorInternalServerError()
	}

	return nil
}
