package repository

import (
	"database/sql"

	"github.com/hegdeshashank73/glamr-backend/common"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/sirupsen/logrus"
)

func GetMagiclink(tx *sql.Tx, arg entities.GetMagiclinkArg) (entities.GetMagiclinkRet, errors.GlamrError) {
	ret := entities.GetMagiclinkRet{}
	query := `SELECT token, email FROM auth_magiclink WHERE token = $1;`
	err := tx.QueryRow(query, arg.Token).Scan(
		&ret.Token,
		&ret.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ret, errors.GlamrErrorGeneralNotFound("magiclink not found")
		}
		logrus.Error(err)
		return ret, errors.GlamrErrorInternalServerError()
	}

	return ret, nil
}

func DeleteMagiclink(tx *sql.Tx, arg entities.DeleteMagiclinkArg) errors.GlamrError {
	query := `DELETE FROM auth_magiclink WHERE token = $1;`
	if _, err := tx.Exec(query, arg.Token); err != nil {
		logrus.Error(err)
		return errors.GlamrErrorDatabaseIssue()
	}
	return nil
}

func CreateAndGetAuthUser(tx *sql.Tx, arg entities.CreateAndGetAuthUserArg) (entities.CreateAndGetAuthUserRet, errors.GlamrError) {
	ret := entities.CreateAndGetAuthUserRet{}
	newId := common.GenerateSnowflake()
	query := `INSERT IGNORE INTO auth_users (id, email) VALUES ($1,$2);`
	result, err := tx.Exec(query, newId, arg.Email)
	if err != nil {
		logrus.Error(err)
		return ret, errors.GlamrErrorInternalServerError()
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return ret, errors.GlamrErrorInternalServerError()
	}
	ret.IsNewUser = rowsAffected > 0

	if ret.IsNewUser {
		ret.AuthUser.Id = newId
		ret.AuthUser.Email = arg.Email
		return ret, nil
	}

	query = `SELECT id, email FROM auth_users WHERE email = $1;`
	err = tx.QueryRow(query, arg.Email).Scan(
		&ret.AuthUser.Id,
		&ret.AuthUser.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return ret, errors.GlamrErrorGeneralNotFound("user not found")
		}
		logrus.Error(err)
		return ret, errors.GlamrErrorInternalServerError()
	}
	return ret, nil
}

func CreateAccessToken(tx *sql.Tx, arg entities.CreateAccessTokenArg) (entities.CreateAccessTokenRet, errors.GlamrError) {
	ret := entities.CreateAccessTokenRet{}
	query := `INSERT INTO auth_tokens (token, user_id) VALUES ($1,$2);`
	_, err := tx.Exec(query, arg.AccessToken, arg.Id)
	if err != nil {
		logrus.Error(err)
		return ret, errors.GlamrErrorInternalServerError()
	}

	ret.AccessToken = arg.AccessToken
	return ret, nil
}

func UpdateNotifToken(tx *sql.Tx, person entities.Person, arg entities.UpdateNotifTokenArg) errors.GlamrError {
	query := `INSERT INTO notifs_tokens (user_id, device_id, fcm_token) VALUES ($1,$2,$3) ON DUPLICATE KEY UPDATE fcm_token = $3;`
	_, err := tx.Exec(query, person.Id, arg.DeviceID, arg.FCMToken)
	if err != nil {
		logrus.Error(err)
		return errors.GlamrErrorInternalServerError()
	}

	return nil
}

func CreateMagiclink(tx *sql.Tx, arg entities.CreateMagiclinkArg) errors.GlamrError {
	query := `INSERT INTO auth_magiclink (token, email) VALUES ($1, $2);`
	if _, err := tx.Exec(query, arg.Token, arg.Email); err != nil {
		logrus.Error(err)
		return errors.GlamrErrorDatabaseIssue()
	}
	return nil
}
