package entities

import "github.com/hegdeshashank73/glamr-backend/common"

type CreatePersonArg struct {
	Id        common.Snowflake
	Username  string
	FirstName string
	LastName  string
	Bio       string
}
