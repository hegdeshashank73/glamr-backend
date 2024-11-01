package repository

import (
	"time"

	"github.com/hegdeshashank73/glamr-backend/dsql"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/utils"
)

func GetPerson(tx dsql.Tx, person *entities.Person, arg entities.GetPersonArg) (entities.GetPersonRet, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("repository.GetPerson", st)

	ret := entities.GetPersonRet{}

	// To implement Get Person

	return ret, nil
}
