package entities

import (
	"github.com/hegdeshashank73/glamr-backend/common"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/utils"
)

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

type CreateMagiclinkReq struct {
	Email string `json:"email"`
}

func (r *CreateMagiclinkReq) Validate() errors.GlamrError {
	email, err := utils.AssertString(r.Email, utils.StringOpts{IsRequired: true, MaxLength: 128, EntityName: "email", ConvertToLowercase: true})
	if err != nil {
		return errors.GlamrErrorGeneralBadRequest(err.Error())
	}
	r.Email = email
	return nil
}

type CreateMagiclinkRes struct {
	Message Message `json:"message"`
}

type VerifyMagiclinkReq struct {
	MagicToken string `json:"token"`
	FCMToken   string `json:"fcm_token"`
	DeviceID   string `json:"device_id"`
	OAuthUser  OUser
}

func (r *VerifyMagiclinkReq) Validate() errors.GlamrError {
	token, err := utils.AssertString(r.MagicToken, utils.StringOpts{IsRequired: true, MaxLength: 128, EntityName: "token"})
	if err != nil {
		return errors.GlamrErrorGeneralBadRequest(err.Error())
	}
	r.MagicToken = token
	return nil
}

type VerifyMagiclinkRes struct {
	Message     Message `json:"message"`
	AccessToken string  `json:"access_token"`
}

type Message struct {
	Img         string `json:"img"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type OUser struct {
	Email        string
	FirstName    string
	LastName     string
	ProfileImage string
}

type GetMagiclinkRet struct {
	Token string
	Email string
}

type GetMagiclinkArg struct {
	Token string
}
type DeleteMagiclinkArg struct {
	Token string
}

type CreateAndGetAuthUserArg struct {
	Email string
}

type CreateAndGetAuthUserRet struct {
	AuthUser  AuthUser
	IsNewUser bool
}
type AuthUser struct {
	Id    common.Snowflake
	Email string
}

type CreateAccessTokenArg struct {
	Id          common.Snowflake
	AccessToken string
}

type CreateAccessTokenRet struct {
	AccessToken string
}
type UpdateNotifTokenArg struct {
	DeviceID string
	FCMToken string
}

type CreateMagiclinkArg struct {
	Token string
	Email string
}
