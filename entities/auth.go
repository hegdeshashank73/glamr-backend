package entities

import "github.com/hegdeshashank73/glamr-backend/common"

type Person struct {
	Id    common.Snowflake `json:"id"`
	About About            `json:"about"`
}

type About struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Img       string `json:"img"`
}

type GetPersonReq struct {
	Token  string           `json:"token"`
	UserID common.Snowflake `json:"user_id"`
}

type GetPersonArg GetPersonReq

type GetPersonRes struct {
	Person Person `json:"person"`
}

type GetPersonRet GetPersonRes
