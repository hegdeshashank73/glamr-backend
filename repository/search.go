package repository

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/hegdeshashank73/glamr-backend/common"
	"github.com/hegdeshashank73/glamr-backend/dsql"
	"github.com/hegdeshashank73/glamr-backend/entities"
	"github.com/hegdeshashank73/glamr-backend/errors"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/sirupsen/logrus"
)

func CreatePersonSearch(tx dsql.Tx, person entities.Person, arg entities.SearchOptionsArg, serpAPIResults map[string]interface{}) (entities.CreatePersonSearchRet, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("repository.CreateSearchOptions", st)

	ret := entities.CreatePersonSearchRet{}

	apiResponse, err := json.Marshal(serpAPIResults)
	if err != nil {
		logrus.Error("failed to marshal json,", err)
		return ret, errors.GlamrErrorInternalServerError()
	}

	id := common.GenerateSnowflake()
	query := `INSERT INTO people_searches (id, user_id, s3_key, country_code, api_response, created_at) VALUES (?,?,?,?,?,?);`
	_, err = tx.Exec(query, id, person.Id, arg.S3Key, arg.CountryCode, apiResponse, time.Now().Unix())
	if err != nil {
		logrus.Error("failed to insert into people_searches,", err)
		return ret, errors.GlamrErrorInternalServerError()
	}

	ret.ID = id
	ret.APIResponse = apiResponse
	return ret, nil
}

func CreateSearchOptions(tx *sql.Tx, arg entities.CreateSearchOptionsArg) errors.GlamrError {
	st := time.Now()
	defer utils.LogTimeTaken("repository.CreateSearchOptions", st)

	query := `INSERT INTO searches_options (id, search_id, title, link, source, source_icon, in_stock, price, image, display_order, currency) VALUES `
	values := []interface{}{}
	for i, option := range arg.SearchOptions {
		query += `(?,?,?,?,?,?,?,?,?,?,?),`
		values = append(values, common.GenerateSnowflake(), arg.ID, option.Title, option.Link, option.Source, option.SourceIcon, option.InStock, option.Price, option.Image, i, option.Currency)
	}
	query = strings.TrimSuffix(query, ",")

	stmt, err := tx.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return errors.GlamrErrorDatabaseIssue()
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		logrus.Error(err)
		return errors.GlamrErrorDatabaseIssue()
	}

	return nil
}

func GetSearchHistory(tx *sql.Tx, person entities.Person) (entities.GetSearchHistoryRet, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("repository.GetSearchHistory", st)

	ret := entities.GetSearchHistoryRet{}

	query := `SELECT id, s3_key, created_at FROM people_searches WHERE user_id = ? ORDER BY created_at DESC;`
	rows, err := tx.Query(query, person.Id)
	if err != nil {
		logrus.Error(err)
		return ret, errors.GlamrErrorDatabaseIssue()
	}
	defer rows.Close()

	searchHistoryList := []entities.SearchHistory{}
	for rows.Next() {
		var searchHistory entities.SearchHistory

		if err := rows.Scan(&searchHistory.ID, &searchHistory.S3Key, &searchHistory.CreatedAt); err != nil {
			logrus.Error(err)
			return ret, errors.GlamrErrorDatabaseIssue()
		}
		searchHistoryList = append(searchHistoryList, searchHistory)
	}
	ret.SearchHistory = searchHistoryList
	return ret, nil
}

func GetSearchHistoryOptions(tx *sql.Tx, arg entities.GetSearchHistoryOptionsArg) (entities.GetSearchHistoryOptionsRes, errors.GlamrError) {
	st := time.Now()
	defer utils.LogTimeTaken("repository.GetSearchHistoryOptions", st)
	ret := entities.GetSearchHistoryOptionsRes{}

	query := `SELECT title, link, source, source_icon, in_stock, price, image, currency FROM searches_options WHERE search_id = ? ORDER BY display_order;`
	rows, err := tx.Query(query, arg.SearchID)
	if err != nil {
		logrus.Error(err)
		return ret, errors.GlamrErrorDatabaseIssue()
	}

	searchOptions := []entities.SearchOptions{}
	for rows.Next() {
		var searchOption entities.SearchOptions
		if err := rows.Scan(&searchOption.Title, &searchOption.Link, &searchOption.Source, &searchOption.SourceIcon, &searchOption.InStock, &searchOption.Price, &searchOption.Image, &searchOption.Currency); err != nil {
			logrus.Error(err)
			return ret, errors.GlamrErrorDatabaseIssue()
		}
		searchOptions = append(searchOptions, searchOption)
	}
	ret.SearchOptions = searchOptions

	return ret, nil
}
