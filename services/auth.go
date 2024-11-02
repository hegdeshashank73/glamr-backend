package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/repository"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/hegdeshashank73/glamr-backend/vendors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func VerifyMagiclink(req entities.VerifyMagiclinkReq) (entities.VerifyMagiclinkRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("services.VerifyMagiclink", st)

	res := entities.VerifyMagiclinkRes{}
	tx, err := vendors.DBMono.Begin()
	if err != nil {
		return res, errors.GlamrErrorDatabaseIssue()
	}

	var magiclink entities.GetMagiclinkRet
	var derr errors.GlamrError

	if req.MagicToken != "" {
		magiclink, derr = repository.GetMagiclink(tx, entities.GetMagiclinkArg{Token: req.MagicToken})
		if derr != nil {
			tx.Rollback()
			return res, derr
		}

		derr = repository.DeleteMagiclink(tx, entities.DeleteMagiclinkArg{Token: req.MagicToken})
		if derr != nil {
			tx.Rollback()
			return res, derr
		}
	}

	firstName := utils.StringOrDefault(req.OAuthUser.FirstName, "")
	lastName := utils.StringOrDefault(req.OAuthUser.LastName, "")
	email := utils.StringOrDefault(magiclink.Email, req.OAuthUser.Email)
	if email == "" {
		return res, errors.GlamrErrorGeneralBadRequest("invalid user")
	}

	ret, derr := repository.CreateAndGetAuthUser(tx, entities.CreateAndGetAuthUserArg{
		Email: email,
	})
	if derr != nil {
		tx.Rollback()
		return res, derr
	}

	userId := ret.AuthUser.Id

	if ret.IsNewUser && req.OAuthUser.ProfileImage == "" {
		_, err = vendors.S3Client.CopyObject(&s3.CopyObjectInput{
			Bucket:     aws.String(viper.GetString("USER_ASSETS_S3_BUCKET")),
			Key:        aws.String(fmt.Sprintf("/%s/avatar.jpg", userId.String())),
			CopySource: aws.String(viper.GetString("ASSETS_S3_BUCKET") + "/defaults/user.jpg"),
		})
		if err != nil {
			logrus.Error(err)
			tx.Rollback()
			return res, errors.GlamrErrorInternalServerError()
		}
	} else if ret.IsNewUser && req.OAuthUser.ProfileImage != "" {
		response, err := http.Get(req.OAuthUser.ProfileImage)
		if err != nil {
			_, err = vendors.S3Client.CopyObject(&s3.CopyObjectInput{
				Bucket:     aws.String(viper.GetString("USER_ASSETS_S3_BUCKET")),
				Key:        aws.String(fmt.Sprintf("/%s/avatar.jpg", userId.String())),
				CopySource: aws.String(viper.GetString("ASSETS_S3_BUCKET") + "/defaults/user.jpg"),
			})
			if err != nil {
				logrus.Error(err)
				tx.Rollback()
				return res, errors.GlamrErrorInternalServerError()
			}
		}
		defer response.Body.Close()

		imageData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logrus.Error(err)
			tx.Rollback()
			return res, errors.GlamrErrorInternalServerError()
		}
		imageReader := bytes.NewReader(imageData)

		_, err = vendors.S3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(viper.GetString("USER_ASSETS_S3_BUCKET")),
			Key:    aws.String(fmt.Sprintf("/%s/avatar.jpg", userId.String())),
			Body:   imageReader,
		})
		if err != nil {
			logrus.Error(err)
			tx.Rollback()
			return res, errors.GlamrErrorInternalServerError()
		}
	}

	if derr = repository.CreatePerson(tx, &entities.CreatePersonArg{
		Id:        userId,
		FirstName: firstName,
		LastName:  lastName,
	}); derr != nil {
		tx.Rollback()
		return res, derr
	}

	atRet, derr := repository.CreateAccessToken(tx, entities.CreateAccessTokenArg{
		Id:          userId,
		AccessToken: utils.GenerateAccessToken(),
	})
	if derr != nil {
		tx.Rollback()
		return res, derr
	}

	if req.FCMToken != "" {
		repository.UpdateNotifToken(tx, entities.Person{Id: userId}, entities.UpdateNotifTokenArg{
			DeviceID: req.DeviceID,
			FCMToken: req.FCMToken,
		})
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorDatabaseIssue()
	}

	res.AccessToken = atRet.AccessToken
	return res, nil
}

func CreateMagiclink(req entities.CreateMagiclinkReq) (entities.CreateMagiclinkRes, errors.GlamrError) {
	res := entities.CreateMagiclinkRes{}
	if !utils.ValidateEmail(req.Email) {
		return res, errors.GlamrErrorInvalidValue("email")
	}

	tx, err := vendors.DBMono.Begin()
	if err != nil {
		return res, errors.GlamrErrorDatabaseIssue()
	}

	token := uuid.New()

	derr := repository.CreateMagiclink(tx, entities.CreateMagiclinkArg{Token: token.String(), Email: req.Email})
	if derr != nil {
		tx.Rollback()
		return res, derr
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return res, errors.GlamrErrorDatabaseIssue()
	}

	if derr = SendMagiclink(req.Email, entities.DataSendMagiclink{
		Magiclink: fmt.Sprintf("%s/verify?token=%s", viper.GetString("BASEURL_WEB"), token),
	}); derr != nil {
		return res, derr
	}

	res.Message = entities.Message{
		Title: "We have emailed you the magiclink 🪄",
	}
	return res, nil
}
