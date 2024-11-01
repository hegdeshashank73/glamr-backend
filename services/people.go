package services

import (
	"time"

	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/repository"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/hegdeshashank73/glamr-backend/vendors"
)

func GetPersonFromToken(token string) (entities.GetPersonRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.GetPersonFromToken", st)

	return GetPerson(nil, entities.GetPersonReq{Token: token})
}

func GetPerson(person *entities.Person, req entities.GetPersonReq) (entities.GetPersonRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.GetPerson", st)

	res := entities.GetPersonRes{}
	ret, derr := repository.GetPerson(vendors.DBMono, person, entities.GetPersonArg(req))
	if derr != nil {
		return res, derr
	}
	res.Person = ret.Person
	return res, nil
}
